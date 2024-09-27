import './style.css'

import * as THREE from 'three';
import fishImage from './fish/image.png'
import { VRButton } from 'three/addons/webxr/VRButton.js';
// import { GLTFLoader } from 'three/addons/loaders/GLTFLoader.js';
import ImmersiveControls from '@depasquale/three-immersive-controls';

const width = window.innerWidth, height = window.innerHeight;

// init

const camera = new THREE.PerspectiveCamera( 70, width / height, 0.01, 10 );
camera.position.z = 1;

const scene = new THREE.Scene();

const geometry = new THREE.BoxGeometry( 0.2, 0.2, 0.2 );
const material = new THREE.MeshNormalMaterial();

// add 2D image
const texture = new THREE.TextureLoader().load(fishImage)
const imageMaterial = new THREE.MeshBasicMaterial({ map: texture, transparent: true })
const imageGeometry = new THREE.PlaneGeometry(0.2, 0.2)
const imageMesh = new THREE.Mesh(imageGeometry, imageMaterial)
scene.add(imageMesh)

const mesh = new THREE.Mesh( geometry, material );
mesh.position.x = 1
mesh.position.y = 2
scene.add( mesh );

const renderer = new THREE.WebGLRenderer( { antialias: true, alpha:true } );
renderer.xr.enabled = true;
renderer.setSize( width, height );
renderer.setAnimationLoop( animate );
document.body.appendChild( renderer.domElement );

const controls = new ImmersiveControls(camera, renderer, scene, { /* options */ });


const light = new THREE.AmbientLight( 0x404040, 38 ); // soft white light
scene.add( light );

// animation

function animate( time: number ) {

	mesh.rotation.x = time / 2000;
	mesh.rotation.y = time / 1000;


	// imageMesh.rotation.x = time / 2000;
  // oscelate the image
  imageMesh.position.y = 2 + Math.sin(time / 1000) / 10
  imageMesh.rotation.y = Math.sin(time / 1000)

  controls.update();
  renderer.render( scene, camera );

}

document.body.appendChild( VRButton.createButton( renderer ) );


