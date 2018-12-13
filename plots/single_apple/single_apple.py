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

with open('stats.dat','r') as csvfile:
    plots = csv.reader(csvfile, delimiter=',')
    for row in plots:
        rmax.append(int(row[0]))
        rmin.append(int(row[1]))
        ravg.append(float(row[2]))

        gmax.append(int(row[4]))
        gmin.append(int(row[5]))
        gavg.append(float(row[6]))

        bmax.append(int(row[8]))
        bmin.append(int(row[9]))
        bavg.append(float(row[10]))


plt.plot(ravg, 'g-', rmax, 'g--', rmin, 'g-.', color='r')
plt.plot(gavg, 'g-', gmax, 'g--', gmin, 'g-.', color='g')
plt.plot(bavg, 'g-', bmax, 'g--', bmin, 'g-.', color='b')
#plt.plot(rmax)
#plt.plot(rmin)
plt.ylabel('r value')
#xlabel('Item (s)')
#ylabel('Value')
#title('Python Line Chart: Plotting numbers')
#grid(True)
plt.show()
