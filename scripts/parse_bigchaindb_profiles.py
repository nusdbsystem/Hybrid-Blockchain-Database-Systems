import pstats
import sys

if len(sys.argv) < 2:
    print("Usage: {} <profile file>".format(sys.argv[0]))
    sys.exit(-1)


p = pstats.Stats(sys.argv[1])
# p.strip_dirs().sort_stats('cumulative').print_stats(30)
p.sort_stats('cumulative').print_stats(30)
