# Notification Plugin

> **Status:** Stable  
> **Platforms:** Android, iOS  
> **Version:** v1 alpha

The Notification plugin posts and cancels system notifications on the device.

---

## Quick Example

```go
package main

import (
    "github.com/sweet-juice/sweetjuice/plugins/notification"
)

func postNotification() {
    notification.Post(notification.Notification{
        Title:    "New Message",
        Body:     "You have a new message",
        Icon:     "ic_notification",
        ChannelID: "default_channel",
    })
}

func cancelNotification() {
    notification.Cancel(1)
}
```

---

## API Reference

### `notification.NewPlugin() *NotificationPlugin`

Creates a new notification plugin instance.

---

### `(*NotificationPlugin).Post(n Notification) (string, error)`

Posts a system notification.

| Field | Type | Description |
|-------|------|-------------|
| `Title` | `string` | Notification title |
| `Body` | `string` | Notification body text |
| `Icon` | `string` | Icon resource name |
| `ChannelID` | `string` | Notification channel ID |
| `ChannelName` | `string` | Notification channel name |

---

### `(*NotificationPlugin).Cancel(id int) string`

Cancels a previously posted notification by ID.

---

## Android Requirements

```xml
<uses-permission android:name="android.permission.POST_NOTIFICATIONS" />
```

On Android 13+, the `POST_NOTIFICATIONS` runtime permission must be granted.

---

## iOS Notes

- iOS requires user permission to show notifications.
- Configure notification categories in your app delegate.
