# PD Notify

Daemon that notifies on PagerDuty events.

## Notifications

This programs sends notifications via OSC escape sequences. It uses OSC 777 by
default, but if your terminal does not support that, you can use the flag
`-use-osc-9` instead, which should be more widely supported.

## Tmux Integration

To use pd-notify with tmux and actually get the desktop notifications, you must
set the `allow-passthrough` setting to allow your terminal to get the correct
escape sequences. For more information, see
https://github.com/tmux/tmux/wiki/FAQ#what-is-the-passthrough-escape-sequence-and-how-do-i-use-it.
An example tmux.conf would look like:

```tmux
set-option -g allow-passthrough
```
