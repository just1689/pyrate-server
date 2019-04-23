class Stash {
    static engine
    static scene
    static canvas
    static camera
    static materials = new Map()
    static ws
    static sceneItems = new Map()
    static light
    static water

    static box
    static mapTiles = new Map()
    static mapTilesBin = new Map()
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
    const sphere = BABYLON.Mesh.CreateSphere("sphere", 16, 10, Stash.scene)
    sphere.position.y = 7
    sphere.material = Stash.materials.get("woodMaterial")

    Stash.box = BABYLON.MeshBuilder.CreateBox("box", {height: 1, width: 16, depth: 16}, Stash.scene)
    Stash.box.material = Stash.materials.get("soilMaterial")
    Stash.box.position.y = 2
    Stash.box.visibility = 0

    Stash.materials.get("waterMaterial").addToRenderList(sphere)
    Stash.materials.get("waterMaterial").addToRenderList(Stash.box)
}


function ConnectWS() {
    Stash.ws = new WebSocket("ws://localhost:8000/ws/test/er");
    Stash.ws.onopen = wsOnOpen
    Stash.ws.onmessage = wsOnMessage
    Stash.ws.onclose = wsOnClose
}

function wsOnOpen() {
    console.log("Socket open...")
}

function handleTile(t) {
    if (Stash.mapTilesBin.size === 0) {
        createTileByClone(t)
    }
}

function createTileByClone(t) {
    let box = Stash.box.clone("box" + t.X + t.Y)
    box.material = Stash.materials.get("soilMaterial")
    box.position.y = 2
    box.position.x = t.X * 16
    box.position.z = t.Y * 16
    box.tag = t
    box.visibility = 1
    Stash.mapTiles.set(t.ID, box)

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