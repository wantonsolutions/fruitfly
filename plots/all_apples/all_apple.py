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

i = "1"
plt.ylim(0,255)
plt.ylabel('RGB value')
plt.ylabel('Photo Session')
with open('stats_all.dat','r') as csvfile:
    plots = csv.reader(csvfile, delimiter=',')
    for row in plots:
        if i != row[0]:

            print(i)
            plt.plot(ravg, 'g-', rmax, 'g--', rmin, 'g-.', color='r')
            plt.plot(gavg, 'g-', gmax, 'g--', gmin, 'g-.', color='g')
            plt.plot(bavg, 'g-', bmax, 'g--', bmin, 'g-.', color='b')
            plt.savefig(i + ".png")
            plt.clf()
            i = row[0]
            del rmax[:]
            del rmin[:]
            del ravg[:]

            del gmax[:]
            del gmin[:]
            del gavg[:]

            del bmax[:]
            del bmin[:]
            del bavg[:]


        rmax.append(int(row[1]))
        rmin.append(int(row[2]))
        ravg.append(float(row[3]))

        gmax.append(int(row[5]))
        gmin.append(int(row[6]))
        gavg.append(float(row[7]))

        bmax.append(int(row[9]))
        bmin.append(int(row[10]))
        bavg.append(float(row[11]))

    plt.plot(ravg, 'g-', rmax, 'g--', rmin, 'g-.', color='r')
    plt.plot(gavg, 'g-', gmax, 'g--', gmin, 'g-.', color='g')
    plt.plot(bavg, 'g-', bmax, 'g--', bmin, 'g-.', color='b')
    plt.ylabel('RGB value')
    plt.savefig(i + ".png")
    plt.clf()

