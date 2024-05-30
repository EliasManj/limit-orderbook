#!/bin/bash

url='http://localhost:8080/order/'
content_type='Content-Type: application/json'

# Default number of requests
num_requests=20

# Function to generate the JSON payload
generate_payload() {
    order_type=$1
    side=$2
    price=$3
    qty=$4

    cat <<EOF
{
    "order_type":"${order_type}",
    "side":"${side}",
    "price":${price},
    "qty":${qty}
}
EOF
}

# Function to print usage
usage() {
    echo "Usage: $0 [-n number_of_requests]"
    exit 1
}

# Parse command-line options
while getopts ":n:" opt; do
    case ${opt} in
        n)
            num_requests=$OPTARG
            ;;
        *)
            usage
            ;;
    esac
done

# Validate number of requests
if ! [[ "$num_requests" =~ ^[0-9]+$ ]]; then
    echo "Error: Number of requests must be a positive integer."
    usage
fi

# Loop to create the specified number of requests
for i in $(seq 1 $num_requests); do
    if (( $i <= num_requests / 2 )); then
        side="Buy"
    else
        side="Sell"
    fi

    price=$(echo "scale=2; 100 + $i * 0.5" | bc)
    qty=$((10 + $i))

    payload=$(generate_payload "GoodTilCancelled" $side $price $qty)

    curl --location "$url" \
         --header "$content_type" \
         --data "$payload"

    echo # Add a new line for better readability in the terminal
done
