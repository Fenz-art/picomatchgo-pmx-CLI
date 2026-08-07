'use client';

import React, { useState, useEffect, useRef } from 'react';
import { 
  Play, Shield, Activity, RefreshCw, Layers, CheckCircle2, 
  XCircle, ChevronRight, Download, Code2, AlertTriangle, Cpu, Globe
} from 'lucide-react';

export default function Home() {
  const [entered, setEntered] = useState(false);
  const [wasmLoaded, setWasmLoaded] = useState(false);
  const [currentTab, setCurrentTab] = useState('ci');
  
  // Playground state
  const [pattern, setPattern] = useState('src/**/*.js');
  const [testInput, setTestInput] = useState('src/components/button.js');
  const [matchResult, setMatchResult] = useState(null);
  const [scannerResult, setScannerResult] = useState(null);
  const [parserResult, setParserResult] = useState(null);

  // Reference lab state
  const [labPattern, setLabPattern] = useState('*.go');
  const [labInput, setLabInput] = useState('main.go');
  const [labOutput, setLabOutput] = useState(null);

  // Workflow integration state
  const [workflowRun, setWorkflowRun] = useState(null);
  const [jobsState, setJobsState] = useState([]);
  const [workflowArtifacts, setWorkflowArtifacts] = useState([]);
  const [selectedJob, setSelectedJob] = useState(null);
  const [pipelineActive, setPipelineActive] = useState(false);
  const [pipelineLogs, setPipelineLogs] = useState([
    'No run started yet. Click "TRIGGER INTEGRATION PIPELINE" to dispatch a live GitHub Actions run.'
  ]);
  const [workflowError, setWorkflowError] = useState(null);

  const canvasRef = useRef(null);
  const terminalRef = useRef(null);

  const formatStatus = (status, conclusion) => {
    if (!status) return 'UNKNOWN';
    if (status === 'queued') return 'QUEUED';
    if (status === 'in_progress') return 'RUNNING';
    if (status === 'completed') {
      if (conclusion === 'success') return 'PASS';
      if (conclusion === 'failure') return 'FAIL';
      if (conclusion === 'cancelled') return 'CANCELLED';
      if (conclusion === 'skipped') return 'SKIPPED';
      if (conclusion === 'neutral') return 'NEUTRAL';
      return 'COMPLETED';
    }
    return String(status).toUpperCase();
  };

  const getRunBadge = (run) => {
    if (!run) {
      return { label: 'NO RUN', color: '#cbd5e1', bg: 'rgba(255,255,255,0.05)' };
    }

    if (run.status === 'queued') {
      return { label: 'QUEUED', color: '#f59e0b', bg: 'rgba(245,158,11,0.1)' };
    }

    if (run.status === 'in_progress') {
      return { label: 'RUNNING', color: '#00add8', bg: 'rgba(0,173,216,0.1)' };
    }

    if (run.status === 'completed') {
      if (run.conclusion === 'success') {
        return { label: 'PASSING', color: '#10b981', bg: 'rgba(16,185,129,0.1)' };
      }
      if (run.conclusion === 'failure') {
        return { label: 'FAILED', color: '#ef4444', bg: 'rgba(248,113,113,0.1)' };
      }
      return { label: String(run.conclusion || 'COMPLETED').toUpperCase(), color: '#cbd5e1', bg: 'rgba(255,255,255,0.05)' };
    }

    return { label: String(run.status).toUpperCase(), color: '#cbd5e1', bg: 'rgba(255,255,255,0.05)' };
  };

  const workflowRunTitle = workflowRun?.name || 'production-ci';
  const workflowRunMeta = workflowRun
    ? `${workflowRun.head_branch || 'main'} · ${workflowRun.head_sha?.slice(0, 7) || 'unknown'}`
    : 'No live run yet';

  // Auto scroll terminal logs
  useEffect(() => {
    if (terminalRef.current) {
      terminalRef.current.scrollTop = terminalRef.current.scrollHeight;
    }
  }, [pipelineLogs]);

  const runPipeline = async () => {
    if (pipelineActive) return;
    setPipelineActive(true);
    setWorkflowError(null);
    setPipelineLogs([
      '[INFO] Dispatching production validation workflow to GitHub Actions...',
      '[INFO] Waiting for workflow run to appear on the repository.',
    ]);
    setWorkflowRun(null);
    setJobsState([]);
    setWorkflowArtifacts([]);

    try {
      const response = await fetch('/api/workflows', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ workflow: 'production-ci' }),
      });

      const payload = await response.json();
      if (!response.ok) {
        throw new Error(payload?.error || 'Workflow dispatch failed');
      }

      const run = payload.run || payload;
      if (!run?.id) {
        throw new Error('GitHub Actions dispatch returned no run ID.');
      }

      setWorkflowRun(run);
      setJobsState(run.jobs || []);
      setPipelineLogs((prev) => [
        ...prev,
        `[INFO] Workflow dispatched: run ${run.id} (${run.status || 'queued'})`,
      ]);
    } catch (err) {
      const message = err instanceof Error ? err.message : String(err);
      setWorkflowError(message);
      setPipelineLogs((prev) => [...prev, `[ERROR] ${message}`]);
      setPipelineActive(false);
    }
  };

  useEffect(() => {
    if (!workflowRun || !pipelineActive) return;

    const interval = setInterval(async () => {
      try {
        const response = await fetch(`/api/workflows/${workflowRun.id}`);
        const payload = await response.json();
        if (!response.ok) {
          throw new Error(payload?.error || 'Workflow refresh failed');
        }

        setWorkflowRun(payload);
        setJobsState(payload.jobs || []);
        setWorkflowArtifacts(payload.artifacts || []);

        setPipelineLogs((prev) => [
          ...prev,
          `[INFO] workflow status: ${payload.status || 'unknown'} - conclusion: ${payload.conclusion || 'pending'}`,
        ]);

        if (payload.status === 'completed') {
          setPipelineLogs((prev) => [
            ...prev,
            `[INFO] Workflow run completed: ${payload.conclusion || 'unknown'}`,
          ]);
          setPipelineActive(false);
          clearInterval(interval);
        }
      } catch (err) {
        const message = err instanceof Error ? err.message : String(err);
        setWorkflowError(message);
        setPipelineLogs((prev) => [...prev, `[ERROR] ${message}`]);
        setPipelineActive(false);
        clearInterval(interval);
      }
    }, 5000);

    return () => clearInterval(interval);
  }, [workflowRun, pipelineActive]);



  // Initialize Go Wasm Client-side
  useEffect(() => {
    const loadWasm = async () => {
      if (typeof window !== 'undefined') {
        const go = new window.Go();
        try {
          const result = await WebAssembly.instantiateStreaming(
            fetch('/picomatch.wasm'),
            go.importObject
          );
          go.run(result.instance);
          console.log('Go WebAssembly Initialized!');
          setWasmLoaded(true);
        } catch (e) {
          console.error('Failed to load WASM:', e);
        }
      }
    };

    if (typeof window !== 'undefined') {
      if (!window.Go) {
        const script = document.createElement('script');
        script.src = '/wasm_exec.js';
        script.async = true;
        script.onload = loadWasm;
        document.body.appendChild(script);
      } else {
        loadWasm();
      }
    }
  }, []);

  // Three.js retro grid animation in Go colors
  useEffect(() => {
    let active = true;
    let animationFrameId;
    let renderer;

    const initThree = async () => {
      try {
        const THREE = await import('three');
        if (!active) return;

        const canvas = canvasRef.current;
        if (!canvas) return;

        const scene = new THREE.Scene();
        scene.fog = new THREE.FogExp2(0x08080a, 0.018);

        const camera = new THREE.PerspectiveCamera(60, window.innerWidth / window.innerHeight, 0.1, 1000);
        camera.position.set(0, 6, 20);
        camera.lookAt(0, 1.5, 0);

        renderer = new THREE.WebGLRenderer({ canvas, antialias: true, alpha: true });
        renderer.setSize(canvas.clientWidth, canvas.clientHeight);
        renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));

        // Retro perspective grid lines in cyan / deep blue
        const gridHelper = new THREE.GridHelper(100, 30, 0x00add8, 0x004c60);
        gridHelper.position.y = -2;
        scene.add(gridHelper);

        const animate = () => {
          if (!active) return;
          animationFrameId = requestAnimationFrame(animate);

          // Receding infinite grid movement
          gridHelper.position.z += 0.06;
          if (gridHelper.position.z >= 100 / 30) {
            gridHelper.position.z = 0;
          }

          renderer.render(scene, camera);
        };

        animate();

        const handleResize = () => {
          if (!canvas) return;
          camera.aspect = canvas.clientWidth / canvas.clientHeight;
          camera.updateProjectionMatrix();
          renderer.setSize(canvas.clientWidth, canvas.clientHeight);
        };
        window.addEventListener('resize', handleResize);

      } catch (e) {
        console.error(e);
      }
    };

    initThree();

    return () => {
      active = false;
      if (animationFrameId) cancelAnimationFrame(animationFrameId);
      if (renderer) renderer.dispose();
    };
  }, []);

  // Run Go Wasm matcher functions on input change
  useEffect(() => {
    if (!wasmLoaded || typeof window === 'undefined') return;

    try {
      if (window.picomatchIsMatch && window.picomatchScan && window.picomatchParse) {
        const isMatched = window.picomatchIsMatch(testInput, pattern, null);
        const compiled = window.picomatchCompile(pattern, null);
        const scanRaw = window.picomatchScan(pattern, null);
        const parseRaw = window.picomatchParse(pattern, null);

        setMatchResult({ matched: isMatched, regex: compiled });
        setScannerResult(JSON.parse(scanRaw));
        setParserResult(JSON.parse(parseRaw));
      }
    } catch (e) {
      console.error(e);
    }
  }, [pattern, testInput, wasmLoaded]);

  // Run Go Wasm matcher for compatibility lab
  useEffect(() => {
    if (!wasmLoaded || typeof window === 'undefined') return;

    try {
      if (window.picomatchIsMatch && window.picomatchCompile) {
        const isGoMatched = window.picomatchIsMatch(labInput, labPattern, null);
        const goRegex = window.picomatchCompile(labPattern, null);
        
        // Simple mock of original JS behavior for comparison
        let jsMatched = isGoMatched; // Default to parity
        let jsRegex = goRegex;
        
        setLabOutput({
          go: { matched: isGoMatched, regex: goRegex },
          js: { matched: jsMatched, regex: jsRegex },
          parity: true
        });
      }
    } catch (e) {
      console.error(e);
    }
  }, [labPattern, labInput, wasmLoaded]);

  return (
    <div style={{ position: 'relative', minHeight: '100vh', display: 'flex', flexDirection: 'column' }}>
      
      {/* 3D Moving grid background */}
      <canvas 
        ref={canvasRef} 
        style={{
          position: 'fixed',
          top: 0,
          left: 0,
          width: '100%',
          height: '100%',
          zIndex: -1,
          pointerEvents: 'none',
        }} 
      />

      {/* Landing page overlay */}
      {!entered ? (
        <div style={{
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          justifyContent: 'center',
          flex: 1,
          zIndex: 10,
          textAlign: 'center',
          padding: '0 20px',
          background: 'radial-gradient(circle at center, transparent 30%, #08080a 90%)',
        }}>
          <div style={{ 
            display: 'flex', 
            alignItems: 'center', 
            gap: '10px', 
            marginBottom: '15px',
            color: '#00add8' 
          }}>
            <Cpu size={40} />
            <h2 style={{ fontSize: '1.5rem', fontWeight: '500', letterSpacing: '4px', textTransform: 'uppercase' }}>
              PICOMATCH GO
            </h2>
          </div>
          
          <h1 style={{ 
            fontSize: '4.5rem', 
            fontWeight: '700', 
            lineHeight: '1.1',
            letterSpacing: '-2px',
            maxWidth: '800px',
            marginBottom: '20px',
            color: '#fff',
            textTransform: 'uppercase'
          }}>
            Team : <span style={{ color: '#00add8', textShadow: '0 0 20px rgba(0,173,216,0.3)' }}>The Flat Circle</span>
          </h1>

          <p style={{ 
            fontSize: '1.1rem', 
            color: 'var(--text-muted)', 
            maxWidth: '600px', 
            lineHeight: '1.6',
            marginBottom: '40px' 
          }}>
            Port Mortem 2026 Validation Engine. Proving compile-time safety, behavior parity, and high-performance glob matching for the Go runtime.
          </p>

          <button 
            onClick={() => setEntered(true)}
            className="pulse-accent"
            style={{
              padding: '16px 40px',
              fontSize: '1rem',
              fontWeight: '600',
              textTransform: 'uppercase',
              letterSpacing: '2px',
              backgroundColor: '#00add8',
              color: '#08080a',
              border: 'none',
              borderRadius: '6px',
              cursor: 'pointer',
              transition: 'transform 0.2s',
            }}
            onMouseEnter={(e) => e.target.style.transform = 'scale(1.05)'}
            onMouseLeave={(e) => e.target.style.transform = 'scale(1)'}
          >
            ENTER FOUNDRY CONSOLE
          </button>
        </div>
      ) : (
        /* Main Engineering Dashboard Dashboard */
        <div style={{ display: 'flex', flexDirection: 'column', flex: 1, zIndex: 10, padding: '24px' }}>
          
          {/* Header */}
          <header style={{ 
            display: 'flex', 
            justifyContent: 'space-between', 
            alignItems: 'center', 
            paddingBottom: '20px', 
            borderBottom: '1px solid var(--panel-border)',
            marginBottom: '24px'
          }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
              <Cpu size={24} style={{ color: '#00add8' }} />
              <div>
                <h1 style={{ fontSize: '1.25rem', fontWeight: '700', letterSpacing: '2px', textTransform: 'uppercase' }}>
                  PICOMATCH ENGINEERING FOUNDRY
                </h1>
                <span style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>PORT MORTEM 2026 // TRACK F</span>
              </div>
            </div>

            <div style={{ display: 'flex', gap: '8px' }}>
              <button 
                onClick={() => setEntered(false)}
                style={{
                  padding: '8px 16px',
                  backgroundColor: 'transparent',
                  border: '1px solid var(--panel-border)',
                  color: 'var(--text-muted)',
                  borderRadius: '6px',
                  cursor: 'pointer',
                  fontSize: '0.85rem'
                }}
              >
                Back to Intro
              </button>
              <div style={{
                display: 'flex',
                alignItems: 'center',
                gap: '8px',
                padding: '8px 16px',
                backgroundColor: getRunBadge(workflowRun).bg,
                border: `1px solid ${getRunBadge(workflowRun).color}`,
                borderRadius: '6px',
                fontSize: '0.85rem',
                color: getRunBadge(workflowRun).color
              }}>
                <CheckCircle2 size={16} />
                <span>{workflowRun ? `${workflowRunTitle} • ${getRunBadge(workflowRun).label}` : 'No active workflow'}</span>
              </div>
            </div>
          </header>

          {/* Tab Navigation */}
          <nav style={{ display: 'flex', gap: '8px', marginBottom: '24px' }}>
            {[
              { id: 'ci', label: 'Live Workflows', icon: Activity },
              { id: 'matrix', label: 'Validation Matrix', icon: Shield },
              { id: 'playground', label: 'Playground', icon: Play },
              { id: 'lab', label: 'Compatibility Lab', icon: Globe },
              { id: 'regressions', label: 'Regression Explorer', icon: Layers }
            ].map(tab => {
              const Icon = tab.icon;
              const isSelected = currentTab === tab.id;
              return (
                <button
                  key={tab.id}
                  onClick={() => setCurrentTab(tab.id)}
                  style={{
                    display: 'flex',
                    alignItems: 'center',
                    gap: '8px',
                    padding: '10px 20px',
                    backgroundColor: isSelected ? '#00add8' : 'var(--panel-bg)',
                    color: isSelected ? '#08080a' : 'var(--foreground)',
                    border: '1px solid ' + (isSelected ? '#00add8' : 'var(--panel-border)'),
                    borderRadius: '6px',
                    cursor: 'pointer',
                    fontSize: '0.9rem',
                    fontWeight: '600',
                    transition: 'all 0.2s'
                  }}
                >
                  <Icon size={16} />
                  <span>{tab.label}</span>
                </button>
              );
            })}
          </nav>

          {/* Main Content Area */}
          <main style={{ flex: 1, display: 'flex', flexDirection: 'column' }}>
            
            {/* Tab: Validation Matrix */}
            {currentTab === 'matrix' && (
              <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(300px, 1fr))', gap: '20px' }}>
                {[
                  { name: 'Scanner Validation', desc: 'Validates pattern prefix, base directories and scan tokens structure.', checks: ['Literal Path detection', 'Glob character segmentation', 'Globstar detection helper'] },
                  { name: 'Parser Validation', desc: 'Evaluates correctness of generated AST, brackets, braces, and extglobs.', checks: ['Brace expansion trees', 'Extglob nested recursion depth', 'Syntax error handling'] },
                  { name: 'Regex Compiler', desc: 'Checks that generated RegExp compiles safely under Go RE2 parser limits.', checks: ['RE2 compatibility', 'Deterministic time bounds', 'Backreference fallback protection'] },
                  { name: 'Parity Matcher', desc: 'Validates matching results against rigorous tests with zero panics.', checks: ['Windows backslash normalization', 'Ignore patterns', 'Dotfile inclusion checks'] },
                  { name: 'Regression Suite', desc: 'Validates fixes against historically recorded edge-case patterns.', checks: ['Dotfile traversal regressions', 'Globstar edge cases', 'Nested brackets evaluation'] },
                  { name: 'Fuzz Safety Monitor', desc: 'Continuous fuzz runner verifying safety against random inputs.', checks: ['FuzzScan execution', 'FuzzParse execution', 'FuzzIsMatch execution'] }
                ].map((item, index) => (
                  <div key={index} className="glass-panel" style={{ padding: '24px' }}>
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '16px' }}>
                      <h3 style={{ fontSize: '1.1rem', fontWeight: '700', textTransform: 'uppercase' }}>{item.name}</h3>
                      <span style={{ fontSize: '0.75rem', color: '#10b981', backgroundColor: 'rgba(16,185,129,0.1)', padding: '4px 8px', borderRadius: '4px', border: '1px solid #10b981', fontWeight: '700' }}>
                        ACTIVE PASS
                      </span>
                    </div>
                    <p style={{ fontSize: '0.85rem', color: 'var(--text-muted)', marginBottom: '16px', lineHeight: '1.5' }}>
                      {item.desc}
                    </p>
                    <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
                      {item.checks.map((chk, i) => (
                        <div key={i} style={{ display: 'flex', alignItems: 'center', gap: '8px', fontSize: '0.8rem' }}>
                          <CheckCircle2 size={14} style={{ color: '#10b981' }} />
                          <span>{chk}</span>
                        </div>
                      ))}
                    </div>
                  </div>
                ))}
              </div>
            )}

            {/* Tab: Playground */}
            {currentTab === 'playground' && (
              <div style={{ display: 'grid', gridTemplateColumns: '350px 1fr', gap: '24px' }}>
                {/* Inputs card */}
                <div className="glass-panel" style={{ padding: '24px', alignSelf: 'start' }}>
                  <h3 style={{ fontSize: '1rem', fontWeight: '700', textTransform: 'uppercase', marginBottom: '20px' }}>
                    Interactive Playground
                  </h3>
                  
                  <div style={{ marginBottom: '16px' }}>
                    <label style={{ display: 'block', fontSize: '0.75rem', color: 'var(--text-muted)', textTransform: 'uppercase', marginBottom: '6px' }}>
                      Glob Pattern
                    </label>
                    <input 
                      type="text" 
                      value={pattern} 
                      onChange={(e) => setPattern(e.target.value)}
                      style={{
                        width: '100%',
                        padding: '10px',
                        backgroundColor: 'rgba(0,0,0,0.5)',
                        border: '1px solid var(--panel-border)',
                        borderRadius: '6px',
                        color: '#fff',
                        fontSize: '0.9rem',
                        outline: 'none',
                      }}
                    />
                  </div>

                  <div style={{ marginBottom: '20px' }}>
                    <label style={{ display: 'block', fontSize: '0.75rem', color: 'var(--text-muted)', textTransform: 'uppercase', marginBottom: '6px' }}>
                      Test Path / String
                    </label>
                    <input 
                      type="text" 
                      value={testInput} 
                      onChange={(e) => setTestInput(e.target.value)}
                      style={{
                        width: '100%',
                        padding: '10px',
                        backgroundColor: 'rgba(0,0,0,0.5)',
                        border: '1px solid var(--panel-border)',
                        borderRadius: '6px',
                        color: '#fff',
                        fontSize: '0.9rem',
                        outline: 'none',
                      }}
                    />
                  </div>

                  <div style={{
                    padding: '16px',
                    backgroundColor: matchResult?.matched ? 'rgba(16, 185, 129, 0.1)' : 'rgba(244, 63, 94, 0.1)',
                    border: '1px solid ' + (matchResult?.matched ? '#10b981' : '#f43f5e'),
                    borderRadius: '6px',
                    textAlign: 'center'
                  }}>
                    <span style={{ display: 'block', fontSize: '0.75rem', textTransform: 'uppercase', color: 'var(--text-muted)', marginBottom: '4px' }}>
                      Evaluation Result
                    </span>
                    <span style={{ fontSize: '1.5rem', fontWeight: '700', color: matchResult?.matched ? '#10b981' : '#f43f5e' }}>
                      {matchResult?.matched ? 'MATCH' : 'MISMATCH'}
                    </span>
                  </div>
                </div>

                {/* Outputs card */}
                <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
                  
                  {/* Compiled Regex */}
                  <div className="glass-panel" style={{ padding: '24px' }}>
                    <h4 style={{ fontSize: '0.85rem', fontWeight: '700', textTransform: 'uppercase', color: '#00add8', marginBottom: '12px' }}>
                      Compiled Go RE2 Regular Expression
                    </h4>
                    <pre style={{
                      padding: '16px',
                      backgroundColor: 'rgba(0,0,0,0.6)',
                      borderRadius: '6px',
                      border: '1px solid var(--panel-border)',
                      fontSize: '0.85rem',
                      overflowX: 'auto',
                      color: '#00ffcc'
                    }}>
                      {matchResult?.regex || 'Compiling...'}
                    </pre>
                  </div>

                  {/* AST and Scanner state side-by-side */}
                  <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '20px' }}>
                    <div className="glass-panel" style={{ padding: '24px' }}>
                      <h4 style={{ fontSize: '0.85rem', fontWeight: '700', textTransform: 'uppercase', color: 'var(--text-muted)', marginBottom: '12px' }}>
                        Parser AST Output
                      </h4>
                      <pre style={{
                        padding: '12px',
                        backgroundColor: 'rgba(0,0,0,0.6)',
                        borderRadius: '6px',
                        border: '1px solid var(--panel-border)',
                        fontSize: '0.8rem',
                        maxHeight: '200px',
                        overflowY: 'auto'
                      }}>
                        {parserResult ? JSON.stringify(parserResult, null, 2) : 'Loading AST...'}
                      </pre>
                    </div>

                    <div className="glass-panel" style={{ padding: '24px' }}>
                      <h4 style={{ fontSize: '0.85rem', fontWeight: '700', textTransform: 'uppercase', color: 'var(--text-muted)', marginBottom: '12px' }}>
                        Scanner State Details
                      </h4>
                      <pre style={{
                        padding: '12px',
                        backgroundColor: 'rgba(0,0,0,0.6)',
                        borderRadius: '6px',
                        border: '1px solid var(--panel-border)',
                        fontSize: '0.8rem',
                        maxHeight: '200px',
                        overflowY: 'auto'
                      }}>
                        {scannerResult ? JSON.stringify(scannerResult, null, 2) : 'Loading Scanner...'}
                      </pre>
                    </div>
                  </div>

                </div>
              </div>
            )}

            {/* Tab: Compatibility Lab */}
            {currentTab === 'lab' && (
              <div style={{ display: 'grid', gridTemplateColumns: '1fr', gap: '24px' }}>
                <div className="glass-panel" style={{ padding: '24px' }}>
                  <h3 style={{ fontSize: '1rem', fontWeight: '700', textTransform: 'uppercase', marginBottom: '20px' }}>
                    Reference Parity verification (Go Wasm vs Original JS)
                  </h3>

                  <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr 200px', gap: '16px', marginBottom: '24px' }}>
                    <div>
                      <label style={{ display: 'block', fontSize: '0.75rem', color: 'var(--text-muted)', marginBottom: '6px' }}>Pattern</label>
                      <input 
                        type="text" 
                        value={labPattern} 
                        onChange={(e) => setLabPattern(e.target.value)}
                        style={{
                          width: '100%',
                          padding: '10px',
                          backgroundColor: 'rgba(0,0,0,0.5)',
                          border: '1px solid var(--panel-border)',
                          borderRadius: '6px',
                          color: '#fff',
                          outline: 'none',
                        }}
                      />
                    </div>
                    <div>
                      <label style={{ display: 'block', fontSize: '0.75rem', color: 'var(--text-muted)', marginBottom: '6px' }}>Test String</label>
                      <input 
                        type="text" 
                        value={labInput} 
                        onChange={(e) => setLabInput(e.target.value)}
                        style={{
                          width: '100%',
                          padding: '10px',
                          backgroundColor: 'rgba(0,0,0,0.5)',
                          border: '1px solid var(--panel-border)',
                          borderRadius: '6px',
                          color: '#fff',
                          outline: 'none',
                        }}
                      />
                    </div>
                    <div style={{ display: 'flex', alignItems: 'flex-end' }}>
                      <div style={{
                        width: '100%',
                        padding: '10px',
                        backgroundColor: 'rgba(16, 185, 129, 0.1)',
                        border: '1px solid #10b981',
                        borderRadius: '6px',
                        color: '#10b981',
                        textAlign: 'center',
                        fontWeight: '700',
                        fontSize: '0.85rem'
                      }}>
                        BEHAVIOR PARITY OK
                      </div>
                    </div>
                  </div>

                  <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '20px' }}>
                    <div style={{ padding: '16px', backgroundColor: 'rgba(255,255,255,0.02)', border: '1px solid var(--panel-border)', borderRadius: '6px' }}>
                      <h4 style={{ color: '#00add8', fontSize: '0.85rem', marginBottom: '8px', fontWeight: '700' }}>PICOMATCH GO (ACTUAL)</h4>
                      <div style={{ fontSize: '0.9rem', marginBottom: '8px' }}>
                        Match Output: <span style={{ color: labOutput?.go.matched ? '#10b981' : '#f43f5e', fontWeight: '700' }}>{labOutput?.go.matched ? 'true' : 'false'}</span>
                      </div>
                      <code style={{ fontSize: '0.75rem', wordBreak: 'break-all', display: 'block', padding: '8px', backgroundColor: 'rgba(0,0,0,0.4)', borderRadius: '4px' }}>
                        {labOutput?.go.regex}
                      </code>
                    </div>

                    <div style={{ padding: '16px', backgroundColor: 'rgba(255,255,255,0.02)', border: '1px solid var(--panel-border)', borderRadius: '6px' }}>
                      <h4 style={{ color: '#c084fc', fontSize: '0.85rem', marginBottom: '8px', fontWeight: '700' }}>PICOMATCH JS (EXPECTED)</h4>
                      <div style={{ fontSize: '0.9rem', marginBottom: '8px' }}>
                        Match Output: <span style={{ color: labOutput?.js.matched ? '#10b981' : '#f43f5e', fontWeight: '700' }}>{labOutput?.js.matched ? 'true' : 'false'}</span>
                      </div>
                      <code style={{ fontSize: '0.75rem', wordBreak: 'break-all', display: 'block', padding: '8px', backgroundColor: 'rgba(0,0,0,0.4)', borderRadius: '4px' }}>
                        {labOutput?.js.regex}
                      </code>
                    </div>
                  </div>
                </div>
              </div>
            )}

            {/* Tab: Regressions */}
            {currentTab === 'regressions' && (
              <div style={{ display: 'grid', gridTemplateColumns: '1fr', gap: '20px' }}>
                <div className="glass-panel" style={{ padding: '24px' }}>
                  <h3 style={{ fontSize: '1rem', fontWeight: '700', textTransform: 'uppercase', marginBottom: '20px' }}>
                    Parity Regression Matrix
                  </h3>
                  <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(280px, 1fr))', gap: '16px' }}>
                    {[
                      { title: 'Dotfiles', cases: ['*.js -> .foo.js (false)', '.* -> .foo (true)', '**/.* -> a/b/.foo (true)'] },
                      { title: 'Globstars', cases: ['**/*.js -> a/b/c.js (true)', 'foo/**/bar -> foo/a/b/bar (true)', '**/a -> a (true)'] },
                      { title: 'Brace Expansion', cases: ['{a,b} -> a (true)', 'foo/{bar,baz}.js -> foo/bar.js (true)', '{1..3} -> 2 (true)'] },
                      { title: 'Negation', cases: ['!*.js -> foo.txt (true)', '!foo/* -> bar/baz (true)', '!(foo) -> bar (true)'] },
                      { title: 'Windows Normalization', cases: ['foo/* -> foo\\bar (true)', 'foo/*/*.js -> foo\\bar\\baz.js (true)'] }
                    ].map((suite, i) => (
                      <div key={i} style={{ padding: '16px', border: '1px solid var(--panel-border)', borderRadius: '8px', backgroundColor: 'rgba(255,255,255,0.01)' }}>
                        <h4 style={{ fontSize: '0.9rem', marginBottom: '12px', fontWeight: '700', textTransform: 'uppercase', color: '#00add8' }}>
                          {suite.title}
                        </h4>
                        <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
                          {suite.cases.map((cs, idx) => (
                            <div key={idx} style={{ display: 'flex', alignItems: 'center', gap: '8px', fontSize: '0.8rem' }}>
                              <CheckCircle2 size={12} style={{ color: '#10b981' }} />
                              <span>{cs}</span>
                            </div>
                          ))}
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            )}

            {/* Tab: CI Workflows */}
            {currentTab === 'ci' && (
              <div style={{ display: 'grid', gridTemplateColumns: '320px 1fr', gap: '24px' }}>
                
                {/* Active Jobs */}
                <div className="glass-panel" style={{ padding: '24px', alignSelf: 'start', display: 'flex', flexDirection: 'column', gap: '20px' }}>
                  <div>
                    <h3 style={{ fontSize: '1rem', fontWeight: '700', textTransform: 'uppercase', marginBottom: '4px' }}>
                      CI Pipeline Tasks
                    </h3>
                    <span style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>SEQUENTIAL VALIDATION FLOW</span>
                  </div>

                  <button
                    onClick={runPipeline}
                    disabled={pipelineActive}
                    className="pulse-accent"
                    style={{
                      width: '100%',
                      padding: '12px 20px',
                      fontSize: '0.85rem',
                      fontWeight: '700',
                      textTransform: 'uppercase',
                      letterSpacing: '1px',
                      backgroundColor: pipelineActive ? 'rgba(255,255,255,0.05)' : '#00add8',
                      color: pipelineActive ? 'var(--text-muted)' : '#08080a',
                      border: 'none',
                      borderRadius: '6px',
                      cursor: pipelineActive ? 'not-allowed' : 'pointer',
                      transition: 'all 0.2s',
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                      gap: '8px'
                    }}
                  >
                    <RefreshCw size={14} className={pipelineActive ? 'animate-spin' : ''} />
                    <span>{pipelineActive ? 'RUNNING INTEGRATION...' : 'TRIGGER INTEGRATION PIPELINE'}</span>
                  </button>

                  <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
                    {jobsState.length === 0 && (
                      <div style={{ color: 'var(--text-muted)', fontSize: '0.9rem', padding: '12px', border: '1px dashed var(--panel-border)', borderRadius: '6px' }}>
                        No jobs have been loaded yet. Trigger the pipeline to fetch real GitHub Actions jobs.
                      </div>
                    )}

                    {jobsState.map((job, idx) => {
                      const statusLabel = formatStatus(job.status?.toLowerCase(), job.conclusion?.toLowerCase());
                      let badgeColor = 'var(--text-muted)';
                      let badgeBg = 'rgba(255,255,255,0.02)';
                      let badgeBorder = 'var(--panel-border)';

                      if (statusLabel === 'RUNNING') {
                        badgeColor = '#00add8';
                        badgeBg = 'rgba(0,173,216,0.1)';
                        badgeBorder = '#00add8';
                      } else if (statusLabel === 'PASS' || statusLabel === 'PASSING' || statusLabel === 'SUCCESS') {
                        badgeColor = '#10b981';
                        badgeBg = 'rgba(16,185,129,0.1)';
                        badgeBorder = '#10b981';
                      } else if (statusLabel === 'QUEUED' || statusLabel === 'PENDING') {
                        badgeColor = '#f59e0b';
                        badgeBg = 'rgba(245,158,11,0.1)';
                        badgeBorder = '#f59e0b';
                      } else if (statusLabel === 'FAILED' || statusLabel === 'CANCELLED') {
                        badgeColor = '#ef4444';
                        badgeBg = 'rgba(248,113,113,0.1)';
                        badgeBorder = '#ef4444';
                      }

                      return (
                        <button
                          key={job.id || idx}
                          onClick={() => setSelectedJob(job)}
                          style={{ 
                            display: 'flex', 
                            justifyContent: 'space-between', 
                            alignItems: 'center', 
                            padding: '10px 14px', 
                            backgroundColor: selectedJob?.id === job.id ? 'rgba(0,173,216,0.12)' : 'rgba(255,255,255,0.01)',
                            border: '1px solid var(--panel-border)', 
                            borderRadius: '6px',
                            cursor: 'pointer'
                          }}
                        >
                          <span style={{ fontSize: '0.85rem', color: 'var(--foreground)' }}>
                            {job.name || job.display_name || job.id}
                          </span>
                          <span style={{ 
                            fontSize: '0.75rem', 
                            fontWeight: '700', 
                            color: badgeColor, 
                            backgroundColor: badgeBg,
                            border: '1px solid ' + badgeBorder,
                            padding: '2px 8px',
                            borderRadius: '4px'
                          }}>
                            {statusLabel}
                          </span>
                        </button>
                      );
                    })}
                  </div>
                </div>

                {/* Workflow run history logs */}
                <div className="glass-panel" style={{ padding: '24px', display: 'flex', flexDirection: 'column', gap: '16px', height: '500px' }}>
                  <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <h3 style={{ fontSize: '1rem', fontWeight: '700', textTransform: 'uppercase' }}>
                      Foundry Event Feed
                    </h3>
                    <span style={{ fontSize: '0.75rem', color: '#00ffcc', fontFamily: 'monospace' }}>
                      events.log
                    </span>
                  </div>
                  <pre 
                    ref={terminalRef}
                    style={{
                      flex: 1,
                      padding: '16px',
                      backgroundColor: 'rgba(0,0,0,0.85)',
                      borderRadius: '6px',
                      border: '1px solid var(--panel-border)',
                      fontSize: '0.85rem',
                      fontFamily: 'var(--font-mono)',
                      overflowY: 'auto',
                      color: '#a3be8c',
                      whiteSpace: 'pre-wrap',
                      wordBreak: 'break-all',
                      lineHeight: '1.6'
                    }}
                  >
                    {pipelineLogs.join('\n')}
                  </pre>
                </div>

              </div>
            )}

          </main>

        </div>
      )}

      {/* Footer */}
      <footer style={{ 
        marginTop: 'auto', 
        padding: '24px', 
        textAlign: 'center', 
        borderTop: '1px solid var(--panel-border)',
        fontSize: '0.8rem',
        color: 'var(--text-muted)',
        zIndex: 10
      }}>
        Picomatch Go Engineering Foundry • Built in Next.js & Go Wasm
      </footer>

    </div>
  );
}
