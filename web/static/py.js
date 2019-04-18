function StartPirates() {

    var canvas = document.getElementById("renderCanvas");

    var createScene = function () {
        var scene = new BABYLON.Scene(engine);

        var camera = new BABYLON.ArcRotateCamera("Camera", 3 * Math.PI / 2, Math.PI / 4, 100, BABYLON.Vector3.Zero(), scene);
        camera.attachControl(canvas, true);

        var light = new BABYLON.HemisphericLight("light1", new BABYLON.Vector3(0, 1, 0), scene);

        // Skybox
        var skybox = BABYLON.Mesh.CreateBox("skyBox", 1000.0, scene);
        var skyboxMaterial = new BABYLON.StandardMaterial("skyBox", scene);
        skyboxMaterial.backFaceCulling = false;
        skyboxMaterial.reflectionTexture = new BABYLON.CubeTexture("static/textures/TropicalSunnyDay", scene);
        skyboxMaterial.reflectionTexture.coordinatesMode = BABYLON.Texture.SKYBOX_MODE;
        skyboxMaterial.diffuseColor = new BABYLON.Color3(0, 0, 0);
        skyboxMaterial.specularColor = new BABYLON.Color3(0, 0, 0);
        skyboxMaterial.disableLighting = true;
        skybox.material = skyboxMaterial;

        // Ground
        var groundTexture = new BABYLON.Texture("static/textures/sand.jpg", scene);
        groundTexture.vScale = groundTexture.uScale = 4.0;

        var groundMaterial = new BABYLON.StandardMaterial("groundMaterial", scene);
        groundMaterial.diffuseTexture = groundTexture;

        var ground = BABYLON.Mesh.CreateGround("ground", 512, 512, 32, scene, false);
        ground.position.y = -1;
        ground.material = groundMaterial;

        // Water
        var waterMesh = BABYLON.Mesh.CreateGround("waterMesh", 512, 512, 32, scene, false);
        var water = new BABYLON.WaterMaterial("water", scene, new BABYLON.Vector2(1024, 1024));
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

    var engine = new BABYLON.Engine(canvas, true, {preserveDrawingBuffer: true, stencil: true});
    var scene = createScene();

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
