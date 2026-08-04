leader = i % numBrokers // standard rr

leader(pInd) = (pInd + startIndex) % numBrokers //rr + shift

shift = startIndex + floor(pInd/numBrokers)

offset = 1 + (shift + j) % (n - 1)

replica_j = (leader(p) + offset) % n


startIndex — solves: fairness across topics. provides entrypoint for RR cycles