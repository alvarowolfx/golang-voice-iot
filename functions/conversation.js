const { dialogflow, DialogflowConversation } = require( 'actions-on-google' );
const i18n = require( '@sfeir/actions-on-google-i18n' );

const RobotArm = require( './robotArm' );

const ACTION_WELCOME = 'Default Welcome Intent';
const ACTION_ARM_SET = 'arm.set';
const ACTION_ARM_MOVE = 'arm.move';
const ACTION_GRIP_OPEN = 'grip.open';
const ACTION_GRIP_CLOSE = 'grip.close';

class Conversation {
  /**
   * @param {RobotArm} robotArm
   */
  constructor( robotArm ) {
    this.robotArm = robotArm;

    const app = dialogflow( );
    i18n
      .configure( {
        directory : `${__dirname}/locales`,
        defaultExtension : 'json',
        defaultLocale : 'en-US',
      } )
      .use( app );

    app.intent( ACTION_WELCOME, this.welcomeHandler.bind( this ) );
    app.intent( ACTION_ARM_MOVE, this.armMoveHandler.bind( this ) );
    app.intent( ACTION_ARM_SET, this.armSetHandler.bind( this ) );
    app.intent( ACTION_GRIP_OPEN, this.gripOpenHandler.bind( this ) );
    app.intent( ACTION_GRIP_CLOSE, this.gripCloseHandler.bind( this ) );

    app.catch( ( conv, err ) => this.errorHandler( err, conv ) );
    app.fallback( ( conv ) => {
      conv.ask( conv.__( 'FALLBACK' ) )
    } )

    this.app = app;
  }

  getDialogFlowApp() {
    return this.app;
  }

  /**
   * @param {DialogflowConversation} conv
   */
  armMoveHandler( conv, { direction, angle } ) {
    return this.robotArm
      .connect()
      .then( ( ) => {  
        let finalAngle = angle      
        if ( !finalAngle ) {
          finalAngle = 30
        }
        
        if ( finalAngle < 0 || finalAngle > 180 ) {
          return conv.ask( conv.__( 'ANGLE_ERROR' ) )          
        }

        const parts = direction.split( '-' )     
        const degree = parts[0]
        const signal = parts[1]
        if ( signal === 'negative' ) {
          finalAngle *= -1
        }

        conv.ask( conv.__( 'ARM.MOVE', { finalAngle, degree } ) );
        return this.robotArm.moveAngle( degree, finalAngle );        
      } )
      .catch( err => this.errorHandler( err, conv ) );
  }

  /**
   * @param {DialogflowConversation} conv
   */
  armSetHandler( conv, { servo, angle } ) {
    return this.robotArm
      .connect()
      .then( ( ) => {        
        let finalAngle = angle      
        if ( !finalAngle ) {
          finalAngle = 20
        }        
        
        if ( finalAngle < 0 || finalAngle > 180 ) {
          return conv.ask( conv.__( 'ANGLE_ERROR' ) )           
        }

        conv.ask( conv.__( 'ARM.SET', { servo, finalAngle } ) )
        return this.robotArm.setAngle( servo, finalAngle );   
        
      } )
      .catch( err => this.errorHandler( err, conv ) );
  }

  /**
   * @param {DialogflowConversation} conv
   */
  gripOpenHandler( conv ) {
    return this.robotArm
      .connect()
      .then( ( ) => {        
        conv.ask( conv.__( 'GRIP.OPEN' ) );
        return this.robotArm.openGrip( );
        
      } );
  }

  /**
   * @param {DialogflowConversation} conv
   */
  gripCloseHandler( conv ) {
    return this.robotArm
      .connect()
      .then( ( ) => {                
        conv.ask( conv.__( 'GRIP.CLOSE' ) );
        return this.robotArm.closeGrip( );
        
      } );
  }

  errorHandler( err, conv ) {
    console.log( 'Catch', err );
    conv.ask( conv.__( 'ERROR' ) );    
  }

  /**
   * @param {DialogflowConversation} conv
   */
  welcomeHandler( conv ) {
    conv.ask( conv.__( 'WELCOME' ) );
  }
}

module.exports = Conversation;
