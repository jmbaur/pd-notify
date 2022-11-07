# PD Notify

Daemon that notifies on PagerDuty events.

## Tmux Integration

To use pd-notify with tmux and actually get the desktop notifications, you must
set the `allow-passthrough` setting to allow your terminal to get the correct
escape sequences. For more information, see
https://github.com/tmux/tmux/wiki/FAQ#what-is-the-passthrough-escape-sequence-and-how-do-i-use-it.
An example tmux.conf would look like:

```tmux
set-option -g allow-passthrough
```
