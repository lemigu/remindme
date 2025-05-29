# remindme
Extremely simple reminder system for the CLI.

## Installation

```bash
go install github.com/lemigu/remindme@latest
```

## Usage

```bash
Usage:
	remindme add <your note here>          -- add a task
	remindme list                          -- list pending tasks
	remindme ack <id>                      -- acknowledge an existing task
```

For quicker use, I'd suggest the following aliases:

```bash
alias reminder="remindme add"
alias reminders="remindme list"
alias ackr="remindme ack"
```

And I personally wanted it to always display reminders on my first login of each day. If you'd like this as well, you can add the following to your `.bashrc`/`.zprofile`/etc:

```bash
LAST_RUN_FILE="~/.reminders_last_timestamp"
TODAY=$(date +%F)

if [ ! -f "$LAST_RUN_FILE" ] || [ "$(cat $LAST_RUN_FILE)" != "$TODAY" ]; then
    remindme list
    echo "$TODAY" > "$LAST_RUN_FILE"
fi
```

