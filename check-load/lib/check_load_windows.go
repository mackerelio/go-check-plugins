//go:build windows

package checkload

// Do the plugin
func Do() {
	panic("The check-load does not work on Windows environment.")
}
