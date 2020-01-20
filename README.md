# QueueReader

This is a simple polling daemon, which, given an SQS queue name and some polling configuration - like number of concurrent polling threads, long polling duration, spawns the desired number of pollers to consume, process and delete messages from the queue.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

1. Install Go - https://golang.org/doc/install

2. Setup AWS credentials in your shell environment -


```
// Linux or OS X
$ export AWS_ACCESS_KEY_ID=YOUR_AKID
$ export AWS_SECRET_ACCESS_KEY=YOUR_SECRET_KEY

// Windows
> set AWS_ACCESS_KEY_ID=YOUR_AKID
> set AWS_SECRET_ACCESS_KEY=YOUR_SECRET_KEY
```

3. Create an SQS queue - https://docs.aws.amazon.com/cli/latest/reference/sqs/create-queue.html

## Running the service

Clone the git repo and navigate to the parent directory

```
$ git clone https://github.com/akhani18/QueueReader.git
...
$ cd QueueReader
```

Start the Poller

```
$ go run main.go -n <SQS-Queue-Name> -l <long-poll-duration> -p <num-pollers> -r <aws-region> 
```
