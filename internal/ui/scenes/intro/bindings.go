package intro

import "github.com/bunnyholes/pokerhole/client/internal/ui/help"

// Bindings returns key binding descriptions for the intro scene.
// Currently, any key skips the intro.
func Bindings() []help.KeyBinding {
	return []help.KeyBinding{
		{Key: "any key", Description: "스킵하여 메인 메뉴로"},
	}
}
