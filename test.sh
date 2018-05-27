#!/bin/bash
set -eu

die() {
		echo "$@"
		exit 1
}

ADDR=127.0.0.1:9999
URL=http://$ADDR/metrics

BASE=$(realpath $(dirname "$0"))

mkdir -p "$BASE/logs"
rm -f "$BASE"/logs/*.log

./multilog_exporter --config.file doc/example.yaml  --metrics.listen-addr $ADDR &
pid=$!
trap 'kill -9 "$pid"' EXIT

sleep 0.5

kernel_panics=$(curl -s "$URL" | sed -rne 's/^kernel_panics ([0-9]*)$/\1/p')
[[ $kernel_panics -ne 0 ]] && die 'kernel_panics != 0'

echo 'foo' > logs/kernel.log
sleep 0.3

kernel_panics=$(curl -s "$URL" | sed -rne 's/^kernel_panics ([0-9]*)$/\1/p')
[[ $kernel_panics -eq 0 ]] || die 'kernel_panics != 0'

echo 'panic: critical problem' >> logs/kernel.log
sleep 0.5
kernel_panics=$(curl -s "$URL" | sed -rne 's/^kernel_panics ([0-9]*)$/\1/p')
[[ $kernel_panics -eq 1 ]] || die "kernel_panics != 1 ($kernel_panics)"


echo "users only in pool test: 30" > logs/instance.log
sleep 0.5
users_only_test=$(curl -s "$URL" | sed -rne 's/^users_only.*test.* ([0-9]*)$/\1/p')
[[ $users_only_test -eq 30 ]] || die "users_only test != 30 ($users_only_test)"

echo "users only in pool prod: 500" >> logs/instance.log
sleep 0.5
users_only_prod=$(curl -s "$URL" | sed -rne 's/^users_only.*prod.* ([0-9]*)$/\1/p')
[[ $users_only_prod -eq 500 ]] || die "users_only prod != 500 ($users_only_prod)"

echo "users only in pool test: 28" >> logs/instance.log
sleep 0.5
users_only_test=$(curl -s "$URL" | sed -rne 's/^users_only.*test.* ([0-9]*)$/\1/p')
[[ $users_only_test -eq 28 ]] || die "users_only test != 28 ($users_only_test)"


echo
echo "All tests finished successfully"
