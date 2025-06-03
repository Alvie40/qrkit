# 🎥 QRKit Video Calling - Manual Test Guide

## Overview
This guide walks you through testing the complete video calling workflow in the QRKit application.

## Prerequisites
- ✅ Server running on http://localhost:8080
- ✅ LiveKit server running on http://localhost:7880
- ✅ Modern browser with WebRTC support (Chrome/Firefox recommended)
- ✅ Camera and microphone permissions enabled

## Test Scenarios

### 1. 🏠 Basic Application Access
**Objective**: Verify all pages load correctly

**Steps**:
1. Open http://localhost:8080/
2. Click on "Employee Interface" link
3. Click on "Debug Interface" link
4. Click on "Comprehensive Test" link
5. Click on "Status Dashboard" link

**Expected Results**:
- All pages load without errors
- Navigation works correctly
- UI elements display properly

---

### 2. 👷 Employee Session Creation
**Objective**: Test session creation and QR code generation

**Steps**:
1. Go to http://localhost:8080/empregado
2. Click "Create New Session" button
3. Observe the generated QR code
4. Note the session ID and room name
5. Copy the client invitation URL

**Expected Results**:
- Session creates successfully
- QR code generates and displays
- Client URL is provided
- Session details are visible

---

### 3. 🔗 Client Connection
**Objective**: Test client joining video session

**Steps**:
1. Copy the client URL from previous test
2. Open in new browser tab/window
3. Allow camera and microphone permissions
4. Observe video connection status

**Expected Results**:
- Client page loads with session ID
- Camera/microphone permissions requested
- Video connection established
- Video streams visible

---

### 4. 🎬 Two-Way Video Communication
**Objective**: Test bidirectional video and audio

**Setup**: Use two browser windows/tabs

**Steps**:
1. **Employee window**: Go to http://localhost:8080/empregado
2. **Employee window**: Create new session
3. **Employee window**: Join the video session
4. **Client window**: Use the generated client URL
5. **Client window**: Join the video session
6. Test audio and video in both directions

**Expected Results**:
- Both participants see each other's video
- Audio works bidirectionally
- No significant latency
- Connection remains stable

---

### 5. 📱 QR Code Mobile Test
**Objective**: Test QR code scanning from mobile device

**Steps**:
1. Create session on desktop browser
2. Use mobile device to scan QR code
3. Allow camera/microphone on mobile
4. Test video call between desktop and mobile

**Expected Results**:
- QR code scans correctly on mobile
- Mobile client connects successfully
- Video works between desktop and mobile
- Audio quality is acceptable

---

### 6. 🔧 Debug and Diagnostics
**Objective**: Use debugging tools to verify system health

**Steps**:
1. Open http://localhost:8080/comprehensive-test
2. Run all automated tests
3. Check WebRTC connectivity
4. Verify LiveKit integration
5. Monitor real-time status at http://localhost:8080/status

**Expected Results**:
- All automated tests pass
- WebRTC ICE candidates generated
- LiveKit server responds
- No critical errors in status dashboard

---

### 7. 🚨 Error Scenarios
**Objective**: Test error handling and recovery

**Test Cases**:
- **No camera/microphone**: Deny permissions, verify graceful handling
- **Network issues**: Simulate poor connection, test reconnection
- **Invalid session**: Try accessing non-existent session ID
- **Multiple clients**: Test multiple clients joining same session

**Expected Results**:
- Appropriate error messages displayed
- Application doesn't crash
- Graceful degradation when possible
- Clear user feedback on issues

---

## Performance Benchmarks

### 🎯 Acceptance Criteria
- **Page Load Time**: < 3 seconds
- **Session Creation**: < 2 seconds
- **Video Connection**: < 5 seconds
- **Audio/Video Sync**: < 100ms latency
- **QR Code Generation**: < 1 second

### 📊 Monitoring Points
- CPU usage during video calls
- Memory consumption
- Network bandwidth utilization
- WebRTC connection statistics
- LiveKit server performance

---

## Troubleshooting Guide

### 🚫 Common Issues

**1. Camera/Microphone Not Working**
- Check browser permissions
- Verify hardware availability
- Test with other applications

**2. Video Connection Fails**
- Check LiveKit server status
- Verify network connectivity
- Review browser console for errors

**3. QR Code Not Generating**
- Check session creation API
- Verify server response
- Review image generation library

**4. Poor Video Quality**
- Check network bandwidth
- Verify WebRTC settings
- Monitor server resources

### 🔍 Debug Resources
- Browser Developer Tools (F12)
- LiveKit server logs
- Application status dashboard
- Network analysis tools

---

## Test Completion Checklist

- [ ] All basic pages accessible
- [ ] Session creation works
- [ ] QR code generation successful
- [ ] Client connection established
- [ ] Two-way video communication
- [ ] Audio transmission verified
- [ ] Mobile QR scanning tested
- [ ] Debug tools validated
- [ ] Error scenarios handled
- [ ] Performance acceptable

---

## Next Steps After Testing

1. **Document Issues**: Record any bugs or performance issues
2. **Optimize Performance**: Address bottlenecks identified
3. **Enhance UI/UX**: Improve user experience based on testing
4. **Security Review**: Validate security measures
5. **Production Deployment**: Prepare for production environment

---

*Last Updated: June 2, 2025*
*QRKit Video Calling Application - Manual Test Guide*
