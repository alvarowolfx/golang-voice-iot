const functions = require( 'firebase-functions' );

const CloudIotCoreDevice = require( './devices/cloudiotcore' );
const RobotArm = require( './robotArm' );
const Conversation = require( './conversation' );

const DEVICE_ID = 'orangepizero';
const PROJECT_ID = process.env.GCLOUD_PROJECT;
const REGION = 'us-central1';
const REGISTRY = 'robot-arm-registry';

const device = new CloudIotCoreDevice( REGION, PROJECT_ID, REGISTRY, DEVICE_ID );
const robotArm = new RobotArm( device );

const conversation = new Conversation( robotArm );
const app = conversation.getDialogFlowApp();

/**
 * Configure Dialogflow Webhook
 */
exports.dialogflowFirebaseFulfillment = functions.https.onRequest( app );
