import matplotlib.pyplot as plt
import csv


rmax = []
rmin = []
ravg = []

gmax = []
gmin = []
gavg = []

bmax = []
bmin = []
bavg = []

plt.title('RGB vales of Averaged Apple')


plt.ylim(0,255)
plt.ylabel('RGB value')
plt.ylabel('Photo Session')

i = "1"
experiments = 0
averageCounter = 0

with open('stats_all.dat','r') as csvfile:
    plots = csv.reader(csvfile, delimiter=',')
    for row in plots:
        if i != row[0]:
            averageCounter = 0
            print(i)
            i = row[0]
        if i == "1":
            rmax.append(int(row[1]))
            rmin.append(int(row[2]))
            ravg.append(float(row[3]))

            gmax.append(int(row[5]))
            gmin.append(int(row[6]))
            gavg.append(float(row[7]))

            bmax.append(int(row[9]))
            bmin.append(int(row[10]))
            bavg.append(float(row[11]))
        else:
            rmax[averageCounter]+=int(row[1])
            rmin[averageCounter]+=int(row[2])
            ravg[averageCounter]+=float(row[3])

            gmax[averageCounter]+=int(row[5])
            gmin[averageCounter]+=int(row[6])
            gavg[averageCounter]+=float(row[7])

            bmax[averageCounter]+=int(row[9])
            bmin[averageCounter]+=int(row[10])
            bavg[averageCounter]+=float(row[11])
        
        averageCounter = averageCounter + 1
    total = int(i)
    for i in range(0,averageCounter):
            rmax[i] = rmax[i]/total
            rmin[i] = rmin[i]/total
            ravg[i] = ravg[i]/total

            gmax[i] = gmax[i]/total
            gmin[i] = gmin[i]/total
            gavg[i] = gavg[i]/total

            bmax[i] = bmax[i]/total
            bmin[i] = bmin[i]/total
            bavg[i] = bavg[i]/total
        
    plt.plot(ravg, 'g-', rmax, 'g--', rmin, 'g-.', color='r')
    plt.plot(gavg, 'g-', gmax, 'g--', gmin, 'g-.', color='g')
    plt.plot(bavg, 'g-', bmax, 'g--', bmin, 'g-.', color='b')
    plt.ylabel('RGB value')
    plt.show()
    #plt.savefig("average_apple.png")
    plt.clf()

