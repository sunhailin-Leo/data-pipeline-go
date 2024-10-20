#!/bin/bash
# Default topic name
TOPIC="data-pipeline-events"
# Default number of generators
NUM=100  # 默认生成 100 条数据
# Docker is not started by default
START=false

# Parsing Input Parameters
while [[ "$#" -gt 0 ]]; do
    case $1 in
        -t|--topic) TOPIC="$2"; shift ;;
        -n|--num) NUM="$2"; shift ;;
        --start) if [[ "$2" == "true" ]]; then START=true; fi; shift ;;
        *) echo "Usage: $0 [-t topic_name] [-n number_of_messages] [--start true]"; exit 1 ;;
    esac
    shift
done

# Check If need to start a Docker container
if [ "$START" == true ] || [ -z "$(docker ps -q -f name=kafka)" ]; then
    echo "Starting Kafka Docker container..."
    docker run -d --name=kafka \
        -p 9092:9092 \
        -e LANG=C.UTF-8 \
        -e KAFKA_INTER_BROKER_LISTENER_NAME=PLAINTEXT \
        -e KAFKA_CONTROLLER_LISTENER_NAMES=CONTROLLER \
        -e CLUSTER_ID=the-custom-id \
        -e KAFKA_NODE_ID=1 \
        -e KAFKA_PROCESS_ROLES=broker,controller \
        -e KAFKA_CONTROLLER_QUORUM_VOTERS="1@127.0.0.1:9093" \
        -e KAFKA_LISTENERS="PLAINTEXT://:9092,CONTROLLER://:9093" \
        -e KAFKA_ADVERTISED_LISTENERS="PLAINTEXT://127.0.0.1:9092" \
        apache/kafka:3.8.0

    # Wait time
    sleep 10
fi

# Check if the Kafka container is running
if [ "$(docker ps -q -f name=kafka)" ]; then
    docker exec -i kafka bash -c "
        cd /opt/kafka

        # Get a list of all existing topics
        TOPICS=\$(bin/kafka-topics.sh --bootstrap-server localhost:9092 --list)

        # Check if the topic exists
        if echo \"\$TOPICS\" | grep -q \"^$TOPIC\$\"; then
            echo 'Topic \"$TOPIC\" already exists. Skipping creation.'
        else
            # Crete the topic
            bin/kafka-topics.sh --create --topic \"$TOPIC\" --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1
            echo 'Topic \"$TOPIC\" created.'
        fi

        # Describe the topic
        bin/kafka-topics.sh --describe --topic \"$TOPIC\" --bootstrap-server localhost:9092

        for i in \$(seq 1 ${NUM}); do

            RANDOM_DATA=\"\$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 10 | head -n 1)\"

            # Send data to topic
            echo \"\${RANDOM_DATA}\" | bin/kafka-console-producer.sh --topic \"$TOPIC\" --bootstrap-server localhost:9092
            echo \"Sent: \${RANDOM_DATA}\"

            # Wait time
            sleep 0.5
        done
    "
else
    echo "Kafka container is not running. Please start the container or check the name."
    exit 1
fi
