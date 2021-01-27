---
layout: docs
title: Cron Scheduler · Docs · MSM
tab: docs
docsitem: concepts-cron
---

Cron Scheduler
==============

From wikipedia:

> Cron is the time-based job scheduler in Unix-like computer operating systems. cron enables users to schedule jobs (commands or shell scripts) to run periodically at certain times or dates.

MSM uses cron to run periodical backups, syncronise worlds that are in RAM back to disk every half hour to prevent data loss in the event of a crash and roll Minecraft log files to prevent them getting too large each night.

See the <a href="https://github.com/msmhq/msm/blob/latest/cron/msm">MSM cron file</a> in the repository for a detailed example.
