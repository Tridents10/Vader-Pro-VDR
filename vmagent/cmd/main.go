package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/Tridents10/Vader-Pro-VDR/vmagent/internal/collector"
	"github.com/Tridents10/Vader-Pro-VDR/vmagent/internal/identity"
	"github.com/Tridents10/Vader-Pro-VDR/vmagent/internal/logger"
	"github.com/Tridents10/Vader-Pro-VDR/vmagent/internal/models"
	"github.com/Tridents10/Vader-Pro-VDR/vmagent/internal/storage"
	"github.com/Tridents10/Vader-Pro-VDR/vmagent/internal/transport"
)

const (
	agentVersion = "1.1.0"
	outputDir    = "output"
)

func main() {
	log, err := logger.New(outputDir)
	if err != nil {
		// If logging cannot start, there is nothing meaningful to do.
		panic(err)
	}
	defer log.Close()

	log.Info("=======================================")
	log.Info("VDR Edge Agent Initialized (Spooler)")
	log.Info("=======================================")

	// 1. Fire off the Initial Baseline Scan immediately upon startup
	performScan(log)

	// 2. Start the Scanner Goroutine (Runs every 7 days)
	// The 'go' keyword runs this loop completely in the background
	go func() {
		// Note: You can change this to time.Minute * 2 for fast testing
		scanTicker := time.NewTicker(7 * 24 * time.Hour)
		defer scanTicker.Stop()

		for range scanTicker.C {
			log.Info("Executing scheduled 7-day scan...")
			performScan(log)
		}
	}()

	// 3. Start the Shipper Loop in the main thread (Retries every 15 minutes)
	// Note: You can change this to time.Second * 30 for fast testing
	shipTicker := time.NewTicker(10 * time.Second)
	defer shipTicker.Stop()

	for range shipTicker.C {
		encPath := filepath.Join(outputDir, "inventory.enc")

		// Check if the encrypted payload exists on the disk
		if _, err := os.Stat(encPath); err == nil {
			log.Info("Found unsent telemetry. Attempting to ship...")

			// Grab the agent UUID to send in the headers
			uuid, idErr := identity.GetOrCreateUUID(outputDir)
			if idErr != nil {
				log.Error("Could not retrieve UUID for shipping: " + idErr.Error())
				continue
			}

			// Read the encrypted file bytes from the disk
			encryptedData, readErr := os.ReadFile(encPath)
			if readErr != nil {
				log.Error("Could not read spool file: " + readErr.Error())
				continue
			}

			// Try to send the payload to the server
			sendErr := transport.SendPayload(uuid, encryptedData)

			if sendErr == nil {
				// Success! The server got it. Delete the local file.
				log.Info("Telemetry shipped successfully. Clearing local spool.")
				os.Remove(encPath)
			} else {
				// Failed! Internet is probably down. Leave the file alone.
				log.Warn("Network unreachable or server error. Will retry in 15 minutes. Error: " + sendErr.Error())
			}
		}
	}
}

// performScan contains all of your original top-to-bottom collection logic
func performScan(log *logger.Logger) {
	log.Info("--- Scan Cycle Started ---")

	// ---- Identity ----
	uuid, err := identity.GetOrCreateUUID(outputDir)
	if err != nil {
		log.Error("UUID Management Failed: " + err.Error())
	} else {
		log.Info("Agent UUID: " + uuid)
	}

	inv := &models.Inventory{
		AgentUUID:    uuid,
		AgentVersion: agentVersion,
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
		Disks:        []models.Disk{},
		Network:      []models.NetworkInterface{},
		Applications: []models.Application{},
	}

	// ---- Host ----
	if host, err := collector.CollectHost(); err != nil {
		log.Error("Host Collection Failed: " + err.Error())
	} else {
		inv.Host = host
	}

	// ---- OS ----
	if osInfo, err := collector.CollectOS(); err != nil {
		log.Error("OS Collection Failed: " + err.Error())
	} else {
		inv.OS = osInfo
	}

	// ---- CPU ----
	if cpuInfo, err := collector.CollectCPU(); err != nil {
		log.Error("CPU Collection Failed: " + err.Error())
	} else {
		inv.CPU = cpuInfo
	}

	// ---- Memory ----
	if memInfo, err := collector.CollectMemory(); err != nil {
		log.Error("Memory Collection Failed: " + err.Error())
	} else {
		inv.Memory = memInfo
	}

	// ---- Disks ----
	if disks, err := collector.CollectDisks(); err != nil {
		log.Error("Disk Collection Failed: " + err.Error())
	} else {
		inv.Disks = disks
	}

	// ---- Network ----
	if net, err := collector.CollectNetwork(); err != nil {
		log.Error("Network Collection Failed: " + err.Error())
	} else {
		inv.Network = net
	}

	// ---- Software ----
	if apps, err := collector.CollectSoftware(); err != nil {
		log.Error("Software Collection Failed: " + err.Error())
	} else {
		inv.Applications = apps
		log.Info("System & Software Profiling Complete")
	}

	// ---- Export ----
	// This will now use the AES-encrypted WriteInventory logic
	if err := storage.WriteInventory(outputDir, inv); err != nil {
		log.Error("Inventory Export Failed: " + err.Error())
	} else {
		log.Info("Encrypted Telemetry Stored Successfully")
	}

	log.Info("--- Scan Cycle Completed ---")
}
