import { Vector3, TextureLoader, MeshBasicMaterial, PlaneGeometry, Mesh, DoubleSide, Group, Texture, ShaderMaterial } from 'three';
import { TextGeometry } from 'three/addons/geometries/TextGeometry.js';
import { FontLoader } from 'three/addons/loaders/FontLoader.js';
import { Game } from './Game';
import font from './helvetiker_regular.typeface.json'
import fishImage from './fish/image.png'
import fragmentShader from './glsl/fragment.glsl?raw'
import vertexShader from './glsl/vertex.glsl?raw'

const defaultTexture = new TextureLoader().load(fishImage)
const fontLoader = new FontLoader();
const parsedFont = fontLoader.parse(font);

const SPEED = 0.009; //how fast the boids travel
const AVOIDANCE_RADIUS = 0.0025; //the radius of the boid's sightline to the walls
const SEP_WEIGHT = 1; //how much the boid separates itself from it's neighboids
const AVO_WEIGHT = 0.2; //how much the boid dodges the walls
const RAN_WEIGHT = 0.003; //how much the boid goes in a random direction
const INERTIA = 0.02; //the proportion with which the rules should affect the current speed


const flock = ['A', 'B'] as const;

export class Boid {
  position: Vector3;
  velocity: Vector3;
  neighborhood: Boid[] = [];
  game: Game
  group: Group
  mesh: Mesh
  material: ShaderMaterial
  flock: typeof flock[number] = flock[Math.floor(Math.random() * flock.length)];

  constructor(game: Game, position: Vector3, velocity: Vector3, name: string = 'fish', texture: Texture = defaultTexture) {
    this.game = game;
    this.position = position;
    this.velocity = velocity;

    // const material = new MeshBasicMaterial({ map: texture, transparent: true, side: DoubleSide, wireframe: true });
    this.material = new ShaderMaterial({
      vertexShader,
      transparent: true,
      fragmentShader,
      uniforms: {
        uTime: { value: 0.0 },
        uTexture: { value: texture }
      },
      // wireframe: true,
      side: DoubleSide
    });
    const geometry = new PlaneGeometry(0.6, 0.6, 20, 20);
    this.mesh = new Mesh(geometry, this.material)
    this.mesh.rotation.set(0, -Math.PI / 2, 0); 

    this.group = new Group();
    this.group.add(this.mesh);

    // text mesh with position and fish name
    const textGeometry = new TextGeometry(name, {
      font: parsedFont,
      size: 0.1,
      height: 0.01,
    });
    const textMaterial = new MeshBasicMaterial({ color: 0xff0000 });
    const textMesh = new Mesh(textGeometry, textMaterial);
    textMesh.rotation.set(0, -Math.PI / 2, 0); 
    textMesh.position.set(0, 0.25, 0);
    // const axesHelper = new AxesHelper( 5 );
    // this.group.add( axesHelper );

    // this.group.add(textMesh);


    this.game.scene.add(this.group);
  }

  move(deltaTime: number) {

		/* apply all rules*/
		const deltaV = this.separation();
		deltaV.add(this.avoidance());
		deltaV.add(this.randomness());
		deltaV.multiplyScalar(INERTIA);

		/* add rules to current velocity and update position */
		this.velocity.add(deltaV);
		this.velocity.clampLength(SPEED*0.5, SPEED);
    this.velocity.z = 0;

		const scaledVel = this.velocity.clone();
		scaledVel.multiplyScalar((deltaTime * 60) / 1000);
		this.position.add(scaledVel);
		// this.position.clamp(this.game.bounds.min, this.game.bounds.max);
		this.updateShape();
	}

  updateShape () {
    // look dir should only be in the xz plane
    const verlocity = this.velocity.clone();
    verlocity.y = 0;
    verlocity.z = 0;
    verlocity.x = -verlocity.x;
    const lookDir = this.group.position.clone().add(this.velocity.clone())
    this.group.lookAt(lookDir);
    this.group.position.set(this.position.x, this.position.y, this.position.z);

    // wave
    this.material.uniforms.uTime.value = performance.now() / 1000;
  }

  /* avoid all boids in neighborhood */
  separation() {
    const result = new Vector3(0, 0, 0);
    for (const b of this.neighborhood) {
      const dist = b.position.distanceTo(this.position);
      const oppositeDir = this.position.clone();
      oppositeDir.sub(b.position);
      if (dist != 0) oppositeDir.divideScalar(dist);
      result.add(oppositeDir);
    }
    result.clampLength(SEP_WEIGHT, SEP_WEIGHT);
    return result;
  }

  /* move away from walls when boid is close to hitting them */
  avoidance() {
    const result = new Vector3(0, 0, 0);
    if (Math.abs(this.position.x) + AVOIDANCE_RADIUS >= 3.85) {
      result.x = -Math.sign(this.position.x);
    }
    // below floor
    if (this.position.y <= 0.3) {
      result.y = 1
    }
    // above ceiling
    if (this.position.y >= 5) {
      result.y = -1
    }
    if (Math.abs(this.position.z) + AVOIDANCE_RADIUS >= 0.5) {
      result.z = -Math.sign(this.position.z);
    }
    result.clampLength(AVO_WEIGHT, AVO_WEIGHT);
    return result;
  }

  /* shake things up with a little randomness */
	randomness() {
		const result = randomVector(1, 1, 1);
		result.clampLength(RAN_WEIGHT, RAN_WEIGHT);
		return result;
	}
}

export function randomVector(xBound: number, yBound: number, zBound: number) {
  const x = Math.random() * 2 * xBound - xBound;
  const y = Math.random() * 2 * yBound - yBound;
  const z = Math.random() * 2 * zBound - zBound;
  return new Vector3(x, y, z);
}