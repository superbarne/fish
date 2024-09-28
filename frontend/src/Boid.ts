import { Vector3, TextureLoader, MeshBasicMaterial, PlaneGeometry, Mesh } from 'three';
import fishImage from './fish/image.png'
import { Game } from './Game';

const SPEED = 0.006; //how fast the boids travel
const AVOIDANCE_RADIUS = 0.0025; //the radius of the boid's sightline to the walls
const SEP_WEIGHT = 1; //how much the boid separates itself from it's neighboids
const AVO_WEIGHT = 0.2; //how much the boid dodges the walls
const RAN_WEIGHT = 0.01; //how much the boid goes in a random direction
const INERTIA = 0.01; //the proportion with which the rules should affect the current speed

const texture = new TextureLoader().load(fishImage)

export class Boid {
  position: Vector3;
  velocity: Vector3;
  mesh: Mesh;
  neighborhood: Boid[] = [];
  game: Game

  constructor(game: Game, position: Vector3, velocity: Vector3) {
    this.game = game;
    this.position = position;
    this.velocity = velocity;

    const material = new MeshBasicMaterial({ map: texture, transparent: true })
    const geometry = new PlaneGeometry(0.45, 0.45)
    this.mesh = new Mesh(geometry, material)

    this.mesh.rotation.set(0, 0, 0);

    this.game.scene.add(this.mesh);
  }

  move(deltaTime: number) {
		/* create neighborhood */
		this.neighborhood = [];
		// for (const b of boids) {
		// 	if (b == this) continue;
		// 	if (this.pos.distanceToSquared(b.pos) <= NEIGHBORHOOD_RADIUS) this.neighborhood.push(b);
		// }

		/* apply all rules*/
		const deltaV = this.separation();
		deltaV.add(this.avoidance());
		deltaV.add(this.randomness());
		deltaV.multiplyScalar(INERTIA);

		/* add rules to current velocity and update position */
		this.velocity.add(deltaV);
		this.velocity.clampLength(SPEED, SPEED);
		const scaledVel = this.velocity.clone();
		scaledVel.multiplyScalar((deltaTime * 60) / 1000);
		this.position.add(scaledVel);
		// this.position.clamp(this.game.bounds.min, this.game.bounds.max);
		this.updateShape(); //update THREE.js shape
	}

  updateShape () {
    this.mesh.position.set(this.position.x, this.position.y, this.position.z);
    /* look in the direction of travel */
    // const lookDir = this.velocity.clone();
    // lookDir.add(this.position);
    // this.mesh.lookAt(lookDir);

    // look dir should only be in the xz plane
    const lookDir = this.velocity.clone();
    lookDir.add(this.position);
    lookDir.z = 0;
    lookDir.x = 0;
    this.mesh.lookAt(lookDir);
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
    if (Math.abs(this.position.x) + AVOIDANCE_RADIUS >= 2.85) {
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