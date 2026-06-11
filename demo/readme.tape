# ABOUTME: vhs script for the short looping README GIF.
# ABOUTME: One typo, one breath, terminal restored, carry on.

Output demo/paush-demo.gif

Set Shell bash
Set FontSize 20
Set Width 1000
Set Height 600
Set Padding 20
Set Theme "Catppuccin Mocha"
Set TypingSpeed 70ms

Hide
Type "source demo/demo-env.sh && export PAUSH_DURATION=4 && paush install > /dev/null && source /tmp/demohome/.zshrc && clear"
Enter
Show

Sleep 1s
Type@45ms "gti status"
Sleep 300ms
Enter
Sleep 5.6s

Sleep 800ms
Type "git status"
Sleep 300ms
Enter
Sleep 3s
