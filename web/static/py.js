class Stash {

    //Constants
    static TILE_SIZE = 16

    //Useful state
    static engine
    static scene
    static canvas
    static camera
    static materials = new Map()
    static ws
    static sceneItems = new Map()
    static light
    static water

    static tileMesh
    static mapTiles = new Map()
    static mapTilesCounter = 0
    static mapTilesBin = new Map()

    static mapOffsetX = 0
    static mapOffsetY = 0

    //JUNK
    static sphereMesh


}

function StartBabylonEngine() {

    createEngine()
    createScene()
    createCamera()
    createMaterials()
    createLight()
    createSkyBox()
    createGround()


    createWater()


    Stash.engine.runRenderLoop(() => {
        if (Stash.scene) {
            Stash.scene.render()
        }
    })

    window.addEventListener("resize", () => {
        Stash.engine.resize()
    })

}


function createEngine() {
    Stash.canvas = document.getElementById("renderCanvas");
    Stash.engine = new BABYLON.Engine(Stash.canvas, true, {preserveDrawingBuffer: true, stencil: true});

}

function createScene() {
    Stash.scene = new BABYLON.Scene(Stash.engine)

}

function createCamera() {
    Stash.camera = new BABYLON.ArcRotateCamera("Camera", 3 * Math.PI / 2, Math.PI / 4, 100, BABYLON.Vector3.Zero(), Stash.scene)
    Stash.camera.attachControl(Stash.canvas, true)

}

function createMaterials() {

    // Sky
    const skyboxMaterial = new BABYLON.StandardMaterial("skyBox", Stash.scene);
    skyboxMaterial.backFaceCulling = false;
    skyboxMaterial.reflectionTexture = new BABYLON.CubeTexture("static/textures/TropicalSunnyDay", Stash.scene);
    skyboxMaterial.reflectionTexture.coordinatesMode = BABYLON.Texture.SKYBOX_MODE;
    skyboxMaterial.diffuseColor = new BABYLON.Color3(0, 0, 0);
    skyboxMaterial.specularColor = new BABYLON.Color3(0, 0, 0);
    skyboxMaterial.disableLighting = true;
    Stash.materials.set("skyboxMaterial", skyboxMaterial)

    // Water
    const waterMaterial = new BABYLON.WaterMaterial("waterMaterial", Stash.scene, new BABYLON.Vector2(1024, 1024))
    waterMaterial.backFaceCulling = true
    waterMaterial.bumpTexture = new BABYLON.Texture("static/textures/waterbump.png", Stash.scene)
    waterMaterial.windForce = -5
    waterMaterial.waveHeight = 0.5
    waterMaterial.bumpHeight = 0.1
    waterMaterial.waveLength = 0.1
    waterMaterial.colorBlendFactor = 0
    Stash.materials.set("waterMaterial", waterMaterial)


    // Ground
    const soilMaterial = new BABYLON.StandardMaterial("soilMaterial", Stash.scene);
    soilMaterial.diffuseTexture = new BABYLON.Texture("static/textures/soil.jpg", Stash.scene);
    Stash.materials.set("soilMaterial", soilMaterial)

    // Items on the map
    const woodMaterial = new BABYLON.StandardMaterial("woodMaterial", Stash.scene);
    woodMaterial.diffuseTexture = new BABYLON.Texture("static/textures/wood.jpg", Stash.scene);
    Stash.materials.set("woodMaterial", woodMaterial)


}


function createLight() {
    Stash.light = new BABYLON.HemisphericLight("light1", new BABYLON.Vector3(0, 1, 0), Stash.scene)

}

function createSkyBox() {
    Stash.sceneItems.set("skybox", BABYLON.Mesh.CreateBox("skyBox", 1000.0, Stash.scene))
    Stash.sceneItems.get("skybox").material = Stash.materials.get("skyboxMaterial")

}

function createGround() {
    const groundTexture = new BABYLON.Texture("static/textures/sand2.jpg", Stash.scene)
    groundTexture.vScale = groundTexture.uScale = 20.0
    const groundMaterial = new BABYLON.StandardMaterial("groundMaterial", Stash.scene)
    groundMaterial.diffuseTexture = groundTexture
    Stash.sceneItems.set("ground", BABYLON.Mesh.CreateGround("ground", 512, 512, 32, Stash.scene, false))
    Stash.sceneItems.get("ground").position.y = -1
    Stash.sceneItems.get("ground").material = groundMaterial

}

function createWater() {
    Stash.water = BABYLON.Mesh.CreateGround("waterMesh", 512, 512, 32, Stash.scene, false)
    Stash.materials.get("waterMaterial").addToRenderList(Stash.sceneItems.get("skybox"))
    Stash.materials.get("waterMaterial").addToRenderList(Stash.sceneItems.get("ground"))
    Stash.water.material = Stash.materials.get("waterMaterial")

}


function playground() {
    Stash.sphereMesh = BABYLON.Mesh.CreateSphere("sphere", Stash.TILE_SIZE, 10, Stash.scene)
    Stash.sphereMesh.position.y = 7
    Stash.sphereMesh.material = Stash.materials.get("woodMaterial")

    Stash.tileMesh = BABYLON.MeshBuilder.CreateBox("box", {
        height: 1,
        width: Stash.TILE_SIZE,
        depth: Stash.TILE_SIZE
    }, Stash.scene)
    Stash.tileMesh.material = Stash.materials.get("soilMaterial")
    Stash.tileMesh.position.y = 2
    Stash.tileMesh.visibility = 0

    Stash.materials.get("waterMaterial").addToRenderList(Stash.sphereMesh)
    Stash.materials.get("waterMaterial").addToRenderList(Stash.tileMesh)
}


function ConnectWS() {
    Stash.ws = new WebSocket("ws://localhost:8000/ws/test/er");
    Stash.ws.onopen = wsOnOpen
    Stash.ws.onmessage = wsOnMessage
    Stash.ws.onclose = wsOnClose
}

function wsOnOpen() {
    console.log("Socket open...")

    //???
    requestMap()
}

function handleTile(tObject) {

    if (Stash.mapTilesBin.size === 0) {
        createTileByClone(tObject)
        return
    }

    let tileMesh
    for (let pair of Stash.mapTilesBin) {
        tileMesh = pair[1]
        break
    }

    Stash.mapTilesBin.delete(tileMesh.tag.ID)
    Stash.mapTiles.set(tObject.ID, tileMesh)
    updateTileToTileObject(tileMesh, tObject)


}

function createTileByClone(tObject) {

    let tileMesh = Stash.tileMesh.clone("box" + Stash.mapTilesCounter++)
    tileMesh.material = Stash.materials.get("soilMaterial")
    updateTileToTileObject(tileMesh, tObject)
    Stash.mapTiles.set(tObject.ID, tileMesh)

}

function updateTileToTileObject(tileMesh, tObject) {
    tileMesh.material = Stash.materials.get("soilMaterial")
    tileMesh.position.y = 2

    tileMesh.position.x = (tObject.X + Stash.mapOffsetX) * Stash.TILE_SIZE
    tileMesh.position.z = (tObject.Y + Stash.mapOffsetY) * Stash.TILE_SIZE
    tileMesh.tag = tObject
    tileMesh.visibility = 1

}


//Removes tiles from view (render) if they are within an area
function GarbageCollectTiles(minX, maxX, minY, maxY) {
    let marked = []
    for (let pair of Stash.mapTiles) {
        let tObject = pair[1].tag
        if (tObject.X >= minX && tObject.X <= maxX && tObject.Y >= minY && tObject.Y <= maxY) {
            marked.push(tObject.ID)
        }
    }

    for (let id of marked) {
        let tileMesh = Stash.mapTiles.get(id)
        Stash.mapTiles.delete(id)
        tileMesh.visibility = 0
        Stash.mapTilesBin.set(id, tileMesh)
    }

}

function alignMap() {
    for (let mapRow of Stash.mapTiles) {
        let tileMesh = mapRow[1]
        tileMesh.position.x = (tileMesh.tag.X + Stash.mapOffsetX) * Stash.TILE_SIZE
        tileMesh.position.z = (tileMesh.tag.Y + Stash.mapOffsetY) * Stash.TILE_SIZE
    }
}


function wsOnMessage(evt) {
    const o = JSON.parse(evt.data)
    if (o.topic === "tile") {
        handleTile(o.body)
        return
    }

    console.log("ws> " + evt.data)
}

function wsOnClose() {
    console.log("Socket closed...")
    setTimeout(ConnectWS, 2000)
}

function sendWs(m) {
    console.log("Sending: " + m)
    Stash.ws.send(m)
}

function requestMap() {
    let o = {
        topic: "map-request",
        body: {
            X: 0,
            Y: 0,
        }
    }
    let msg = JSON.stringify(o)
    sendWs(msg)
}