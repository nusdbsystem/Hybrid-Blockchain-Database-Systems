import sys
from bigchaindb_driver import BigchainDB
from bigchaindb_driver.crypto import generate_keypair
import queue, threading, time

if len(sys.argv) != 5:
    print('Usage: python3 bench.py load_file_path run_file_path endpoints nthread')
    sys.exit(1)

alice, bob = generate_keypair(), generate_keypair()
metadata = {alice.public_key: bob.public_key}

loadFile, runFile, urls, threadNum = sys.argv[1], sys.argv[2], sys.argv[3].split(','), int(sys.argv[4])
bdbs = []
for url in urls:
    bdb = BigchainDB(url)
    bdbs.append(bdb)

print("BigchainDB with {} threads and {} servers.".format(threadNum, len(urls)))

def readFile(filepath, outQueue):
    with open(filepath, 'r', encoding='UTF-8') as f:
        line = f.readline()
        num = 0
        while line is not None and line != '':
            if line.startswith('INSERT') == False and line.startswith('READ') == False and line.startswith('UPDATE') == False:
                line = f.readline()
                continue
            outQueue.put(line)
            line = f.readline()
            num = num + 1
            if num == 10000:
                break

def sendTxn(lineQueue, latQueue, driver):
    while lineQueue.empty() == False:
        start = time.time()
        try:
            line = lineQueue.get(block=False, timeout=0)
        except Empty:
            continue
        args = line.split(' ', 3)
        if "INSERT" in line or "UPDATE" in line:
            data = {
                'data': {
                    args[2]: {
                        args[2]: args[3],
                    },
                },
            }
            prepared_creation_tx = driver.transactions.prepare(
                operation='CREATE',
                signers=alice.public_key,
                asset=data,
                metadata=metadata,
            )
            fulfilled_creation_tx = driver.transactions.fulfill(
                prepared_creation_tx, private_keys=alice.private_key)
            sent_creation_tx = driver.transactions.send_async(fulfilled_creation_tx)
        else:
            driver.assets.get(search=args[2])
        end = time.time()
        if latQueue is not None:
            latQueue.put(end-start)

print("Start loading init data...")
loadQueue = queue.Queue(maxsize=100000)
readFile(loadFile, loadQueue)
#tLoadRead = threading.Thread(target=readFile, args=(loadFile, loadQueue,))
#tLoadRead.start()
#time.sleep(5)
num = loadQueue.qsize()
start = time.time()
loadThreadList = []
for i in range(32):
    t = threading.Thread(target=sendTxn, args=(loadQueue, None, bdbs[i%len(bdbs)],))
    loadThreadList.append(t)
    t.start()
#tLoadRead.join()
for t in loadThreadList:
    t.join()
end = time.time()
print("Load throughput {} TPS".format(num/(end - start)))

print("Start running experiments...")
runQueue = queue.Queue(maxsize=100000)
latencyQueue = queue.Queue(maxsize=100000)

#tRunRead = threading.Thread(target=readFile, args=(runFile, runQueue,))
#tRunRead.start()
#time.sleep(5)
readFile(runFile, runQueue)
time.sleep(5)

runThreadList = []
for i in range(threadNum):
    t = threading.Thread(target=sendTxn, args=(runQueue, latencyQueue, bdbs[i%len(bdbs)],))
    runThreadList.append(t)

start = time.time()

for t in runThreadList:
    t.start()
time.sleep(1)
for t in runThreadList:
    t.join()

end = time.time()

#allLatency = []
#def getLatency(latQueue):
lat = 0
num = 0
while latencyQueue.empty() == False:
    ts = latencyQueue.get()
    lat = lat + ts
    num = num + 1

#        allLatency.append(ts)
#tLatency = threading.Thread(target=getLatency, args=(latencyQueue,))
#tLatency.start()

# print("Before join...")
# tRunRead.join()
#for t in runThreadList:
#    t.join()

print('Throughput of {} txn: {} txn/s'.format(num, num/(end-start)))
print('Latency: {} ms'.format(lat/num*1000))
