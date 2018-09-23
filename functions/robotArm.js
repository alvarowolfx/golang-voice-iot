class RobotArm {
  /**
   * @param {device} CloudIoTCoreDevice device
   */
  constructor( device ) {
    this.device = device;    
  }

  connect() {
    return this.device.connect();
  }

  /**
   * @param {String} degree - which degree - elbow | shoulder | base
   * @param {Number} angle - Angle to set the arm
   */
  setAngle( degree, angle ) {    
    const config = { [degree] : String( angle ) };
    return this.device.sendConfigToDevice( config );
  }

  /**
   * @param {String} degree - which degree - elbow | shoulder | base
   * @param {Number} angle - Angle to increment on the arm
   */
  moveAngle( degree, angle ) {    
    const config = { [`move${degree}`] : String( angle ) };
    return this.device.sendConfigToDevice( config );
  }

  openGrip() {
    return this.updateGrip( 'open' );
  }

  closeGrip() {
    return this.updateGrip( 'close' );
  }

  /**
   * @param {string} state - open | close
   */
  updateGrip( state ) {    
    const config = { grip : String( state ) };
    return this.device.sendConfigToDevice( config );
  }

}

module.exports = RobotArm;
