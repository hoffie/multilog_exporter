#!/bin/bash
set -eu

die() {
		echo "$@"
		exit 1
}

ADDR=127.0.0.1:9999
URL=http://$ADDR/metrics

BASE="$(realpath "$(dirname "$0")")"

rm -rf "$BASE/logs/" "$BASE"/temp/*.yaml
mkdir -p "$BASE/logs" "$BASE/temp"

cp doc/example.yaml temp/example.yaml

go build .
./multilog_exporter --config.file temp/example.yaml  --metrics.listen-addr $ADDR --debug &
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

cp doc/example2.yaml temp/example.yaml
kill -1 "$pid"
sleep 0.5

echo "requests per second: 12" >> logs/instance.log
sleep 0.5
requests_per_second=$(curl -s "$URL" | sed -rne 's/^requests_per_second.* ([0-9]+)$/\1/p')
[[ $requests_per_second -eq 12 ]] || die "requests_per_second != 12 ($requests_per_second)"

echo "foo" >> logs/links-should-also-work.log
sleep 0.5
metric_from_a_symlinked_file=$(curl -s "$URL" | sed -rne 's/^metric_from_a_symlinked_file.* ([0-9]+)$/\1/p')
[[ $metric_from_a_symlinked_file -eq 1 ]] || die "metric_from_a_symlinked_file != 1 ($metric_from_a_symlinked_file)"
rm -f logs/links-should-also-work.log
mkdir -p logs/by-date/ logs/by-date2/
touch logs/by-date/2022.log
(cd logs && ln -s by-date/2022.log links-should-also-work.log)
sleep 0.5

echo "foo" >> logs/links-should-also-work.log
sleep 0.8
metric_from_a_symlinked_file=$(curl -s "$URL" | sed -rne 's/^metric_from_a_symlinked_file.* ([0-9]+)$/\1/p')
[[ $metric_from_a_symlinked_file -eq 2 ]] || die "metric_from_a_symlinked_file != 2 ($metric_from_a_symlinked_file)"

(cd logs && ln -sf by-date/2023.log links-should-also-work.log)
echo "foo" >> logs/links-should-also-work.log
sleep 1
metric_from_a_symlinked_file=$(curl -s "$URL" | sed -rne 's/^metric_from_a_symlinked_file.* ([0-9]+)$/\1/p')
[[ $metric_from_a_symlinked_file -eq 3 ]] || die "metric_from_a_symlinked_file != 3 ($metric_from_a_symlinked_file)"

(cd logs && ln -sf by-date2/2023.log links-should-also-work.log)
echo "foo" >> logs/links-should-also-work.log
sleep 1
metric_from_a_symlinked_file=$(curl -s "$URL" | sed -rne 's/^metric_from_a_symlinked_file.* ([0-9]+)$/\1/p')
[[ $metric_from_a_symlinked_file -eq 4 ]] || die "metric_from_a_symlinked_file != 3 ($metric_from_a_symlinked_file)"


echo "All tests finished successfully"
sleep 6000
