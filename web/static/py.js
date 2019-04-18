function StartPirates() {

    const canvas = document.getElementById("renderCanvas");

    const createScene = function () {
        const scene = new BABYLON.Scene(engine);

        const camera = new BABYLON.ArcRotateCamera("Camera", 3 * Math.PI / 2, Math.PI / 4, 100, BABYLON.Vector3.Zero(), scene);
        camera.attachControl(canvas, true);

        const light = new BABYLON.HemisphericLight("light1", new BABYLON.Vector3(0, 1, 0), scene);

        // Skybox
        const skybox = BABYLON.Mesh.CreateBox("skyBox", 1000.0, scene);
        const skyboxMaterial = new BABYLON.StandardMaterial("skyBox", scene);
        skyboxMaterial.backFaceCulling = false;
        skyboxMaterial.reflectionTexture = new BABYLON.CubeTexture("static/textures/TropicalSunnyDay", scene);
        skyboxMaterial.reflectionTexture.coordinatesMode = BABYLON.Texture.SKYBOX_MODE;
        skyboxMaterial.diffuseColor = new BABYLON.Color3(0, 0, 0);
        skyboxMaterial.specularColor = new BABYLON.Color3(0, 0, 0);
        skyboxMaterial.disableLighting = true;
        skybox.material = skyboxMaterial;

        // Ground
        const groundTexture = new BABYLON.Texture("static/textures/sand.jpg", scene);
        groundTexture.vScale = groundTexture.uScale = 4.0;

        const groundMaterial = new BABYLON.StandardMaterial("groundMaterial", scene);
        groundMaterial.diffuseTexture = groundTexture;

        const ground = BABYLON.Mesh.CreateGround("ground", 512, 512, 32, scene, false);
        ground.position.y = -1;
        ground.material = groundMaterial;

        // Water
        const waterMesh = BABYLON.Mesh.CreateGround("waterMesh", 512, 512, 32, scene, false);
        const water = new BABYLON.WaterMaterial("water", scene, new BABYLON.Vector2(1024, 1024));
        water.backFaceCulling = true;
        water.bumpTexture = new BABYLON.Texture("static/textures/waterbump.png", scene);
        water.windForce = -5;
        water.waveHeight = 0.5;
        water.bumpHeight = 0.1;
        water.waveLength = 0.1;
        water.colorBlendFactor = 0;
        water.addToRenderList(skybox);
        water.addToRenderList(ground);
        waterMesh.material = water;

        return scene;
    }

    const engine = new BABYLON.Engine(canvas, true, {preserveDrawingBuffer: true, stencil: true});
    const scene = createScene();

    engine.runRenderLoop(function () {
        if (scene) {
            scene.render();
        }
    });

// Resize
    window.addEventListener("resize", function () {
        engine.resize();
    });

}
