import React, { useState, useEffect } from 'react';
import { QRCodeSVG } from 'qrcode.react';
import { useNavigate } from 'react-router-dom';
import api from '../services/api';
import './Login.css';

function Login() {
  const [qrData, setQrData] = useState(null);
  const [proofRequestId, setProofRequestId] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    const token = localStorage.getItem('sessionToken');
    if (token) {
      navigate('/dashboard');
    }
  }, [navigate]);

  useEffect(() => {
    if (!proofRequestId) return;

    const interval = setInterval(async () => {
      try {
        const response = await api.get(`/login/status/${proofRequestId}`);
        if (response.data.status === 'completed' && response.data.session_token) {
          localStorage.setItem('sessionToken', response.data.session_token);
          clearInterval(interval);
          navigate('/dashboard');
        } else if (response.data.status === 'not_found') {
          clearInterval(interval);
          setError('Login session expired. Please try again.');
          setQrData(null);
          setProofRequestId(null);
        }
      } catch (err) {
        console.error('Error checking login status:', err);
      }
    }, 2000);

    return () => clearInterval(interval);
  }, [proofRequestId, navigate]);

  const handleLogin = async () => {
    setLoading(true);
    setError(null);

    try {
      const response = await api.post('/login');
      setQrData(response.data.qr_data);
      setProofRequestId(response.data.proof_request_id);
    } catch (err) {
      setError('Failed to generate login QR code. Please try again.');
      console.error('Login error:', err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="login-page">
      <div className="login-container">
        <div className="card">
          <h1>SSI Sign-In</h1>
          <p className="subtitle">Sign in using your digital wallet</p>

          {error && <div className="error-message">{error}</div>}

          {!qrData ? (
            <div className="login-content">
              <p>Click the button below to generate a QR code for sign-in.</p>
              <button
                className="btn btn-primary"
                onClick={handleLogin}
                disabled={loading}
              >
                {loading ? 'Generating...' : 'Generate QR Code'}
              </button>
            </div>
          ) : (
            <div className="qr-content">
              <p>Scan this QR code with your wallet app:</p>
              <div className="qr-wrapper">
                <QRCodeSVG value={qrData} size={256} />
              </div>
              <p className="qr-hint">
                Open your wallet app and scan the QR code to sign in
              </p>
              <p className="qr-info">Waiting for authentication...</p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

export default Login;

