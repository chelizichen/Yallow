# PROJECT

## DEV

````sh
wails dev
````

## BUILD

````sh
wails build
````

## PACKAGE

````sh
cd build/bin
create-dmg yallow.dmg Yallow.app
````

## 已知问题

````sh
$ 1. 系统变量在运行后无法拿到
   例如我想使用marscode 打开项目目录，使用 whereis marscode 拿到了该变量，但是 os.Env.append 后 仍然不生效
````
