class Stash {
    static scene
    static camera
    static materials = new Map()
    static ws
}

function createMaterials() {
    const skyboxMaterial = new BABYLON.StandardMaterial("skyBox", Stash.scene);
    skyboxMaterial.backFaceCulling = false;
    skyboxMaterial.reflectionTexture = new BABYLON.CubeTexture("static/textures/TropicalSunnyDay", Stash.scene);
    skyboxMaterial.reflectionTexture.coordinatesMode = BABYLON.Texture.SKYBOX_MODE;
    skyboxMaterial.diffuseColor = new BABYLON.Color3(0, 0, 0);
    skyboxMaterial.specularColor = new BABYLON.Color3(0, 0, 0);
    skyboxMaterial.disableLighting = true;
    Stash.materials.set("skybox", skyboxMaterial)

    const woodMaterial = new BABYLON.StandardMaterial("woodMaterial", Stash.scene);
    woodMaterial.diffuseTexture = new BABYLON.Texture("static/textures/wood.jpg", Stash.scene);
    Stash.materials.set("woodMaterial", woodMaterial)

    const soilMaterial = new BABYLON.StandardMaterial("soilMaterial", Stash.scene);
    soilMaterial.diffuseTexture = new BABYLON.Texture("static/textures/soil.jpg", Stash.scene);
    Stash.materials.set("soilMaterial", soilMaterial)

}


function StartPirates() {

    const canvas = document.getElementById("renderCanvas");

    const createScene = function () {

        Stash.scene = new BABYLON.Scene(engine);

        Stash.camera = new BABYLON.ArcRotateCamera("Camera", 3 * Math.PI / 2, Math.PI / 4, 100, BABYLON.Vector3.Zero(), Stash.scene);
        Stash.camera.attachControl(canvas, true);

        createMaterials()

        const light = new BABYLON.HemisphericLight("light1", new BABYLON.Vector3(0, 1, 0), Stash.scene);

        const skybox = BABYLON.Mesh.CreateBox("skyBox", 1000.0, Stash.scene);
        skybox.material = Stash.materials.get("skyBox")

        const groundTexture = new BABYLON.Texture("static/textures/sand2.jpg", Stash.scene);
        groundTexture.vScale = groundTexture.uScale = 20.0;
        const groundMaterial = new BABYLON.StandardMaterial("groundMaterial", Stash.scene);
        groundMaterial.diffuseTexture = groundTexture;
        const ground = BABYLON.Mesh.CreateGround("ground", 512, 512, 32, Stash.scene, false);
        ground.position.y = -1;
        ground.material = groundMaterial;


        // Sphere
        const sphere = BABYLON.Mesh.CreateSphere("sphere", 16, 10, Stash.scene);
        sphere.position.y = 7;
        sphere.material = Stash.materials.get("woodMaterial");


        const box = BABYLON.MeshBuilder.CreateBox("box", {height: 1, width: 16, depth: 16}, Stash.scene);
        box.material = Stash.materials.get("soilMaterial")
        box.position.y = 2;


        // Water
        const waterMesh = BABYLON.Mesh.CreateGround("waterMesh", 512, 512, 32, Stash.scene, false);
        const water = new BABYLON.WaterMaterial("water", Stash.scene, new BABYLON.Vector2(1024, 1024));
        water.backFaceCulling = true;
        water.bumpTexture = new BABYLON.Texture("static/textures/waterbump.png", Stash.scene);
        water.windForce = -5;
        water.waveHeight = 0.5;
        water.bumpHeight = 0.1;
        water.waveLength = 0.1;
        water.colorBlendFactor = 0;
        water.addToRenderList(skybox);
        water.addToRenderList(ground);
        water.addToRenderList(sphere);
        // water.addToRenderList(box);
        waterMesh.material = water;

    }

    const engine = new BABYLON.Engine(canvas, true, {preserveDrawingBuffer: true, stencil: true});
    createScene();

    engine.runRenderLoop(function () {
        if (Stash.scene) {
            Stash.scene.render();
        }
    });

    window.addEventListener("resize", function () {
        engine.resize();
    });

}

function ConnectWS() {
    Stash.ws = new WebSocket("ws://localhost:8000/ws/test/er");

    Stash.ws.onopen = function () {

        // Web Socket is connected, send data using send()
        alert("Message is sent...");
    };

    Stash.ws.onmessage = wsMessageIn

    Stash.ws.onclose = function () {

        // websocket is closed.
        alert("Connection is closed...");
    };

}

function wsMessageIn(evt) {
    const msg = evt.data;
    alert(">>>" + msg)
}

function send(m) {
    Stash.ws.send(m)
}