import { PerspectiveCamera, Scene, WebGLRenderer, AmbientLight, BoxGeometry, MeshNormalMaterial, Mesh, PlaneGeometry, MeshBasicMaterial, DoubleSide, Color, TextureLoader, Texture } from "three";
import { Boid, randomVector } from "./Boid";
import ImmersiveControls from '@depasquale/three-immersive-controls';
import { gradientShaderMaterial } from "./Gradient";
import GUI from 'lil-gui'; 

const NUM_BOIDS = 0;

const fishTextureMap = new Map<string, Texture>()

export class Game {
  boids: Boid[] = [];
  width: number;
  height: number;
  scene: Scene;
  renderer: WebGLRenderer;
  camera: PerspectiveCamera;
  controls: ImmersiveControls
  gui: GUI;
  backdropColor1 = new Color(0x00078a)
  backdropColor2 = new Color(0x0a8185)

  floorColor = new Color(0x000030)
  
  testMesh: Mesh;
  constructor(width: number, height: number) {
    this.width = width;
    this.height = height;

    // init
    this.camera = new PerspectiveCamera( 70, width / height, 0.09, 20 );
    this.scene = new Scene();
    this.gui = new GUI();

    // generate boids
    for (let i = 0; i < NUM_BOIDS; i++) {
      const position = randomVector(10, 1.5, 0.5)
      position.y = position.y +2
      const boid = new Boid(this, position, randomVector(1, 1, 1));
      this.boids.push(boid);
    }

    // test mesh
    const geometry = new BoxGeometry( 0.2, 0.2, 0.2 );
    const material = new MeshNormalMaterial();
    this.testMesh = new Mesh( geometry, material );
    this.testMesh.position.x = 0
    this.testMesh.position.y = 0
    // this.scene.add(this.testMesh);

    // floor
    const floorGeometry = new PlaneGeometry( 100, 10 );
    const floorMaterial = new MeshBasicMaterial({ side: DoubleSide });
    floorMaterial.color = this.floorColor;
    const floorPlane = new Mesh( floorGeometry, floorMaterial );
    floorPlane.rotation.x = Math.PI / 2;
    this.scene.add(floorPlane);

    // backdrop
    const backdropGeometry = new PlaneGeometry( 20, 6 );
    // const backdropMaterial = new MeshBasicMaterial( {color: 0x030579, side: DoubleSide } );
    const backdropPlane = new Mesh( backdropGeometry, gradientShaderMaterial( this.backdropColor1, this.backdropColor2 ));
    backdropPlane.position.z = -1.5;
    backdropPlane.position.y = 3;
    this.scene.add(backdropPlane);

    // light
    const light = new AmbientLight( 0xffff00, 60 ); // soft white light
    this.scene.add( light );

    // renderer
    this.renderer = new WebGLRenderer({ antialias: true, alpha:true });
    this.renderer.xr.enabled = true;
    this.renderer.setSize(width, height);
    this.renderer.setAnimationLoop((time) => this.animate(time));

    this.controls = new ImmersiveControls(this.camera, this.renderer, this.scene);

    this.controls.camera.position.set(0, 0.5, 0);

    // const axesHelper = new AxesHelper( 5 );
    // this.scene.add( axesHelper );
    
    this.gui.addColor(this, 'backdropColor1')
    this.gui.addColor(this, 'backdropColor2')
    this.gui.addColor(this, 'floorColor')

    document.body.appendChild(this.renderer.domElement);
    document.body.appendChild( this.renderer.domElement );

    // server stuff
    this.initServerConnection();
  }

  async initServerConnection() {
    const evtSource = new EventSource("/aquarium/38d7976d-3c27-4e74-8bfe-a9ec44318d3f/sse");
    evtSource.addEventListener("ping", (event) => {
      console.log('ping', event.data)
    });

    evtSource.addEventListener("fishjoin", async (event) => {
      console.log('fishjoin', event.data)
      const fish = JSON.parse(event.data);
      console.log(fish)

      // load texture
      if(!fishTextureMap.has(fish.id)) {
        const imageReponse = await fetch(`/aquarium/38d7976d-3c27-4e74-8bfe-a9ec44318d3f/fishes/${fish.filename}`)
        const imageBlob = await imageReponse.blob()
        const texture = new TextureLoader().load(URL.createObjectURL(imageBlob));
        fishTextureMap.set(fish.id, texture)
      }

      // loop for 50 times
      for (let i = 0; i < 50; i++) {
        const position = randomVector(5, 1.3, 1.5)
        position.y = position.y +2
        const boid = new Boid(this, position, randomVector(1, 0, 1), fish.name, fishTextureMap.get(fish.id));
        this.boids.push(boid);
      }
      
    });
  }

  lastTime = 0
  animate(time: number) {
    const deltaTime = time - this.lastTime;
    this.lastTime = time;

    for (const b of this.boids) {
      b.move(deltaTime);
    }

    this.testMesh.rotation.x = time / 2000;
    this.testMesh.rotation.y = time / 1000;


    this.controls.update();
    this.renderer.render( this.scene, this.camera );
    
  }
}