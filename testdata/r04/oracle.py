"""Bounded independent Python zoneinfo oracle over the exact official TZif."""
import datetime as d
import hashlib
import json
import pathlib
import sys
import zoneinfo
p = pathlib.Path(sys.argv[1])
raw = p.read_bytes()
if hashlib.sha256(raw).hexdigest() != "d7f2206b3a45989fc9ad63d558922532fa7352280d5f87176bf1db79cb1d1fa9":
    raise SystemExit("official TZif bytes changed")
with p.open("rb") as f:
    zone = zoneinfo.ZoneInfo.from_file(f, key="America/New_York")
def opened(t):
    local = t.astimezone(zone)
    wall = (local.hour, local.minute, local.second, local.microsecond)
    return local.weekday() == 6 and (1,30,0,0) <= wall < (2,30,0,0)
def text(t):
    return t.isoformat().replace("+00:00", "Z")
start = d.datetime(2026,11,1,5,tzinfo=d.timezone.utc)
transitions = []
previous = opened(start - d.timedelta(seconds=1))
# Exactly 10801 ticks, fixed three-hour domain; not a universal TZ oracle.
for n in range(10801):
    t = start + d.timedelta(seconds=n)
    current = opened(t)
    if current != previous:
        transitions.append({"instant":text(t),"open":current})
    previous = current
states = []
for item in transitions:
    t = d.datetime.fromisoformat(item["instant"].replace("Z", "+00:00"))
    for offset in (-1,0,1):
        point = t + d.timedelta(microseconds=offset)
        states.append({"instant":text(point),"open":opened(point)})
print(json.dumps({"transitions":transitions,"states":states},indent=2))
