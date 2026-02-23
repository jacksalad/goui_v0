package main

import (
	"fmt"
	"github.com/jacksalad/goui_v0/component"
	"github.com/jacksalad/goui_v0/layout"
	"github.com/jacksalad/goui_v0/render"
	"github.com/jacksalad/goui_v0/window"
	"runtime"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Windows Syscall for CPU
var (
	modkernel32        = windows.NewLazySystemDLL("kernel32.dll")
	procGetSystemTimes = modkernel32.NewProc("GetSystemTimes")
)

type FILETIME struct {
	DwLowDateTime  uint32
	DwHighDateTime uint32
}

func (ft FILETIME) ToUint64() uint64 {
	return uint64(ft.DwHighDateTime)<<32 | uint64(ft.DwLowDateTime)
}

func GetSystemTimes() (idle, kernel, user uint64, err error) {
	var i, k, u FILETIME
	r1, _, e1 := procGetSystemTimes.Call(
		uintptr(unsafe.Pointer(&i)),
		uintptr(unsafe.Pointer(&k)),
		uintptr(unsafe.Pointer(&u)),
	)
	if r1 == 0 {
		return 0, 0, 0, e1
	}
	return i.ToUint64(), k.ToUint64(), u.ToUint64(), nil
}

func main() {
	win, err := window.NewWindow(window.WindowConfig{
		Title:     "System Monitor",
		Width:     800,
		Height:    500,
		Resizable: false,
	})
	if err != nil {
		panic(err)
	}

	// Fonts
	font := render.NewFont("Segoe UI", 14)
	// titleFont := render.NewFont("Segoe UI", 20)

	// Root Layout
	win.Root.BgColor = 0xFFF5F5F5 // Light gray background
	win.Root.SetLayout(&layout.HBoxLayout{
		Padding: 20,
		Spacing: 20,
	})

	// Left Column (Stats)
	leftCol := component.NewPanel(0, 0, 370, 460)
	leftCol.BgColor = 0xFFF5F5F5
	leftCol.SetLayout(&layout.VBoxLayout{Spacing: 20})
	win.Root.Add(leftCol)

	// Right Column (More Stats or Logs)
	rightCol := component.NewPanel(0, 0, 370, 460)
	rightCol.BgColor = 0xFFF5F5F5
	rightCol.SetLayout(&layout.VBoxLayout{Spacing: 20})
	win.Root.Add(rightCol)

	// CPU Card
	cpuCard := component.NewCard(370, 200, "CPU Usage")
	leftCol.Add(cpuCard)

	cpuText := component.NewButton("0%")
	cpuText.Font = font
	cpuCard.Add(cpuText)

	cpuProgress := component.NewProgressBar(330, 10)
	cpuProgress.Color = 0xFF0078D7 // Blue
	cpuCard.Add(cpuProgress)

	cpuGraph := component.NewLineChart(330, 80)
	cpuGraph.Color = 0xFF0078D7
	cpuGraph.MaxY = 100
	cpuCard.Add(cpuGraph)

	// Memory Card
	memCard := component.NewCard(370, 200, "Memory Usage")
	leftCol.Add(memCard)

	memText := component.NewButton("0 MB")
	memText.Font = font
	memCard.Add(memText)

	memProgress := component.NewProgressBar(330, 10)
	memProgress.Color = 0xFF881288 // Purple
	memCard.Add(memProgress)

	memGraph := component.NewLineChart(330, 80)
	memGraph.Color = 0xFF881288
	memGraph.MaxY = 100 // Percentage
	memCard.Add(memGraph)

	// Network Card (Simulated)
	netCard := component.NewCard(370, 200, "Network Speed")
	rightCol.Add(netCard)

	netText := component.NewButton("Dl: 0 KB/s  Ul: 0 KB/s")
	netText.Font = font
	netCard.Add(netText)

	netGraph := component.NewLineChart(330, 80)
	netGraph.Color = 0xFF00AA00 // Green
	netGraph.MaxY = 1000        // Dynamic?
	netCard.Add(netGraph)

	// Goroutine Card
	goCard := component.NewCard(370, 100, "Goroutines")
	rightCol.Add(goCard)

	goText := component.NewButton("0")
	goText.Font = font
	goCard.Add(goText)

	// Monitor Loop
	go func() {
		var prevIdle, prevKernel, prevUser uint64
		prevIdle, prevKernel, prevUser, _ = GetSystemTimes()

		// Network sim
		simPhase := 0.0

		ticker := time.NewTicker(500 * time.Millisecond)
		for range ticker.C {
			// CPU
			idle, kernel, user, _ := GetSystemTimes()
			deltaIdle := idle - prevIdle
			deltaKernel := kernel - prevKernel
			deltaUser := user - prevUser

			total := deltaKernel + deltaUser
			cpuUsage := 0.0
			if total > 0 {
				cpuUsage = 100.0 * float64(total-deltaIdle) / float64(total)
			}

			prevIdle = idle
			prevKernel = kernel
			prevUser = user

			// Memory (Go Runtime + Sys Estimate)
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			// Total system memory? Hard without syscall.
			// Let's show App Memory for now as "Memory Usage"
			// Or maybe I can use GlobalMemoryStatusEx via syscall...
			// For simplicity, stick to App Memory relative to 1GB scale for graph visual
			memMB := float64(m.Alloc) / 1024 / 1024
			memUsagePercent := (memMB / 100.0) * 100 // Assume 100MB is "a lot" for this demo
			if memUsagePercent > 100 {
				memUsagePercent = 100
			}

			// Network Sim
			simPhase += 0.1
			netSpeed := 500.0 + 300.0*(simPhase-float64(int(simPhase))) // Random-ish

			// Update UI
			cpuText.Text = fmt.Sprintf("Total: %.1f%%", cpuUsage)
			cpuProgress.SetValue(cpuUsage / 100.0)
			cpuGraph.AddPoint(cpuUsage)

			memText.Text = fmt.Sprintf("App Alloc: %.2f MB", memMB)
			memProgress.SetValue(memUsagePercent / 100.0)
			memGraph.AddPoint(memUsagePercent)

			netText.Text = fmt.Sprintf("Down: %.0f KB/s", netSpeed)
			netGraph.AddPoint(netSpeed)

			goText.Text = fmt.Sprintf("Count: %d", runtime.NumGoroutine())

			win.RequestRepaint()
		}
	}()

	win.Show()
	win.Run()
}
