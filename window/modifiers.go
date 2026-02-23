package window

import "goui/event"

func getModifiers() uint32 {
	var mods uint32
	// VK_SHIFT = 0x10
	if state, _, _ := procGetKeyState.Call(0x10); int16(state)&-32768 != 0 {
		mods |= event.ModShift
	}
	// VK_CONTROL = 0x11
	if state, _, _ := procGetKeyState.Call(0x11); int16(state)&-32768 != 0 {
		mods |= event.ModCtrl
	}
	// VK_MENU (Alt) = 0x12
	if state, _, _ := procGetKeyState.Call(0x12); int16(state)&-32768 != 0 {
		mods |= event.ModAlt
	}
	return mods
}
