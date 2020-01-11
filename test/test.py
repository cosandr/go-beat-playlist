from hashlib import sha1
from random import shuffle

song_hash = sha1()

# path = "C:/Program Files (x86)/Steam/steamapps/common/Beat Saber/Beat Saber_Data/CustomLevels/189721 Night Raid with a Dragon/"
# path = "C:/Program Files (x86)/Steam/steamapps/common/Beat Saber/Beat Saber_Data/CustomLevels/3036 (Milk Crown on Sonnetica - Hexagonial)/"
path = "./nightraid/"

addList = ["info.dat", "Hard.dat", "Expert.dat", "ExpertPlus.dat", "Lightshow.dat"]
# addList = ["info.dat", "ExpertPlus.dat"]

for file in addList:
    with open(F"{path}{file}", "rb") as fr:
        song_hash.update(fr.read())

print(song_hash.hexdigest().lower())
