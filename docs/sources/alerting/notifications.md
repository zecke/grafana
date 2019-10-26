+++
title = "Alerting Notifications"
description = "Alerting Notifications Guide"
keywords = ["Grafana", "alerting", "guide", "notifications"]
type = "docs"
[menu.docs]
name = "Notifications"
parent = "alerting"
weight = 2
+++


# Alert Notifications

> Alerting is only available in Grafana v4.0 and above.

When an alert changes state, it sends out notifications. Each alert rule can have
multiple notifications. In order to add a notification to an alert rule you first need
to add and configure a `notification` channel (can be email, PagerDuty or other integration).
This is done from the Notification Channels page.

## Notification Channel Setup

On the Notification Channels page hit the `New Channel` button to go the page where you
can configure and setup a new Notification Channel.

You specify a name and a type, and type specific options. You can also test the notification to make
sure it's setup correctly.

### Default (send on all alerts)

When checked, this option will notify for all alert rules - existing and new.

### Send reminders

> Only available in Grafana v5.3 and above.

{{< docs-imagebox max-width="600px" img="/img/docs/v53/alerting_notification_reminders.png" class="docs-image--right" caption="Alerting notification reminders setup" >}}

When this option is checked additional notifications (reminders) will be sent for triggered alerts. You can specify how often reminders
should be sent using number of seconds (s), minutes (m) or hours (h), for example `30s`, `3m`, `5m` or `1h` etc.

**Important:** Alert reminders are sent after rules are evaluated. Therefore a reminder can never be sent more frequently than a configured [alert rule evaluation interval](/alerting/rules/#name-evaluation-interval).

These examples show how often and when reminders are sent for a triggered alert.

Alert rule evaluation interval | Send reminders every | Reminder sent every (after last alert notification)
---------- | ----------- | -----------
`30s` | `15s` | ~30 seconds
`1m` | `5m` | ~5 minutes
`5m` | `15m` | ~15 minutes
`6m` | `20m` | ~24 minutes
`1h` | `15m` | ~1 hour
`1h` | `2h` | ~2 hours

<div class="clearfix"></div>

### Disable resolve message

When checked, this option will disable resolve message [OK] that is sent when alerting state returns to false.

## Supported Notification Types

Grafana ships with the following set of notification types:

### All supported notifiers

Name | Type | Supports images | Support alert rule tags
-----|------|---------------- | -----------------------
Prometheus Alertmanager | `prometheus-alertmanager` | yes, external only | yes

# Enable images in notifications {#external-image-store}

Grafana can render the panel associated with the alert rule as a PNG image and include that in the notification. Read more about the requirements and how to configure image rendering [here](/administration/image_rendering/).

Most Notification Channels require that this image be publicly accessible (Slack and PagerDuty for example). In order to include images in alert notifications, Grafana can upload the image to an image store. It currently supports
Amazon S3, Webdav, Google Cloud Storage and Azure Blob Storage. So to set that up you need to configure the [external image uploader](/installation/configuration/#external-image-storage) in your grafana-server ini config file.

Be aware that some notifiers requires public access to the image to be able to include it in the notification. So make sure to enable public access to the images. If you're using local image uploader, your Grafana instance need to be accessible by the internet.

Notification services which need public image access are marked as 'external only'.

# Use alert rule tags in notifications {#alert-rule-tags}

> Only available in Grafana v6.3+.

Grafana can include a list of tags (key/value) in the notification.
It's called alert rule tags to contrast with tags parsed from timeseries.
It currently supports only the Prometheus Alertmanager notifier.

 This is an optional feature. You can get notifications without using alert rule tags.

# Configure the link back to Grafana from alert notifications

All alert notifications contain a link back to the triggered alert in the Grafana instance.
This url is based on the [domain](/installation/configuration/#domain) setting in Grafana.
