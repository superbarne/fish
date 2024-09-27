import { PerspectiveCamera, Scene, WebGLRenderer, AmbientLight, BoxGeometry, MeshNormalMaterial, Mesh, PlaneGeometry, MeshBasicMaterial, DoubleSide } from "three";
import { Boid, randomVector } from "./Boid";
import ImmersiveControls from '@depasquale/three-immersive-controls';

const NUM_BOIDS = 100;

export class Game {
  boids: Boid[] = [];
  width: number;
  height: number;
  scene: Scene;
  renderer: WebGLRenderer;
  camera: PerspectiveCamera;
  controls: ImmersiveControls
  
  testMesh: Mesh;
  constructor(width: number, height: number) {
    this.width = width;
    this.height = height;

    // init
    this.camera = new PerspectiveCamera( 70, width / height, 0.09, 20 );
    this.scene = new Scene();

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
    this.scene.add(this.testMesh);

    // floor
    const floorGeometry = new PlaneGeometry( 100, 10 );
    const floorMaterial = new MeshBasicMaterial( {color: 0x1f00d2, side: DoubleSide } );
    const floorPlane = new Mesh( floorGeometry, floorMaterial );
    floorPlane.rotation.x = Math.PI / 2;
    this.scene.add(floorPlane);

    // backdrop
    const backdropGeometry = new PlaneGeometry( 20, 20 );
    const backdropMaterial = new MeshBasicMaterial( {color: 0x030579, side: DoubleSide } );
    const backdropPlane = new Mesh( backdropGeometry, backdropMaterial );
    backdropPlane.position.z = -3;
    this.scene.add(backdropPlane);

    // light
    const light = new AmbientLight( 0x404040, 38 ); // soft white light
    this.scene.add( light );

    // renderer
    this.renderer = new WebGLRenderer( { antialias: true, alpha:true } );
    this.renderer.xr.enabled = true;
    this.renderer.setSize(width, height);
    this.renderer.setAnimationLoop((time) => this.animate(time));

    this.controls = new ImmersiveControls(this.camera, this.renderer, this.scene);

    document.body.appendChild(this.renderer.domElement);
    document.body.appendChild( this.renderer.domElement );
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