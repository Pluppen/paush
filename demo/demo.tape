# ABOUTME: vhs script for the full paush demo video (16:9 MP4).
# ABOUTME: Story: install, work, typo, breathe, configure deeper breaths, carry on.

Output demo/paush-demo.mp4

Set Shell bash
Set FontSize 22
Set Width 1280
Set Height 720
Set Padding 24
Set Theme "Catppuccin Mocha"
Set TypingSpeed 70ms

Hide
Type "source demo/demo-env.sh"
Enter
Show

Sleep 1.5s
Type "# meet paush — a mindfulness break for your terminal"
Enter
Sleep 1.5s
Type "paush install"
Sleep 500ms
Enter
Sleep 3s

Hide
Type "source /tmp/demohome/.zshrc && clear"
Enter
Show

Sleep 1s

Type "git status"
Sleep 300ms
Enter
Sleep 2s
Type@45ms "gti statu"
Sleep 200ms
Type@45ms "s"
Sleep 300ms
Enter
Sleep 3.6s

Sleep 1s
Type "git status"
Sleep 300ms
Enter
Sleep 2.5s

Type "# breathe deeper:"
Enter
Sleep 1s
Type "export PAUSH_DURATION=6"
Sleep 300ms
Enter
Sleep 1.5s
Type@45ms "sl -la"
Sleep 300ms
Enter
Sleep 7.6s

Sleep 1s
Type "ls -la"
Sleep 300ms
Enter
Sleep 2.5s

Hide
Type "clear"
Enter
Show

Sleep 1s
Type "# a typo is an invitation to pause"
Enter
Sleep 2s
Type "# breathe in. breathe out. carry on."
Enter
Sleep 6s
