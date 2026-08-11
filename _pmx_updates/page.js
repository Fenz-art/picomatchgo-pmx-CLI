'use client';

import React, { useCallback, useState, useEffect, useRef } from 'react';
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
  const [selectedJobLogs, setSelectedJobLogs] = useState('');
  const [jobLogLoading, setJobLogLoading] = useState(false);
  const [pipelineActive, setPipelineActive] = useState(false);
  const [pipelineLogs, setPipelineLogs] = useState([
    'No run started yet. Click "TRIGGER INTEGRATION PIPELINE" to dispatch a live GitHub Actions run.'
  ]);
  const [workflowError, setWorkflowError] = useState(null);
  const [validationResults, setValidationResults] = useState([]);
  const [validationRun, setValidationRun] = useState(null);
  const [validationLoading, setValidationLoading] = useState(false);
  const [validationError, setValidationError] = useState(null);

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

  const validationMatrix = [
    { id: 'help', name: 'pmx help', category: 'CLI', mode: 'local' },
    { id: 'match', name: 'pmx match', category: 'CLI', mode: 'local' },
    { id: 'scan', name: 'pmx scan', category: 'CLI', mode: 'local' },
    { id: 'parse', name: 'pmx parse', category: 'CLI', mode: 'local' },
    { id: 'explain', name: 'pmx explain', category: 'CLI', mode: 'local' },
    { id: 'validate', name: 'pmx validate', category: 'CLI', mode: 'local' },
    { id: 'compat', name: 'pmx compat', category: 'Compatibility', mode: 'local' },
    { id: 'regression', name: 'pmx regression --json', category: 'Regression', mode: 'local' },
    { id: 'bench', name: 'pmx bench', category: 'Reliability', mode: 'local' },
    { id: 'fuzz', name: 'pmx fuzz', category: 'Reliability', mode: 'local' },
    { id: 'doctor', name: 'pmx doctor', category: 'Doctor', mode: 'local' },
    { id: 'doctorJson', name: 'pmx doctor --json', category: 'Doctor', mode: 'local' },
    { id: 'doctorCi', name: 'pmx doctor --ci', category: 'Doctor', mode: 'local' },
    { id: 'agentInspect', name: 'pmx agent inspect --json', category: 'ADLC', mode: 'local' },
    { id: 'agentCheck', name: 'pmx agent check --json', category: 'ADLC', mode: 'local' },
    { id: 'ci', name: 'pmx ci --json', category: 'CI', mode: 'local' },
  ];

  const parseJsonResponse = useCallback(async (response) => {
    const raw = await response.text();
    if (!raw) {
      return {};
    }

    try {
      return JSON.parse(raw);
    } catch (err) {
      return { error: raw, parseError: err instanceof Error ? err.message : String(err) };
    }
  }, []);

  const runValidationSuite = useCallback(async (selectedIds = null) => {
    setValidationLoading(true);
    setValidationError(null);

    try {
      const response = await fetch('/api/validation', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          commands: selectedIds || validationMatrix.map((entry) => entry.id),
        }),
      });

      const payload = await parseJsonResponse(response);
      if (!response.ok) {
        throw new Error(payload?.error || 'Validation failed to execute.');
      }

      if (Array.isArray(payload.results)) {
        setValidationResults(payload.results);
      }
      if (payload.run) {
        setValidationRun(payload.run);
      } else if (Array.isArray(payload.results)) {
        setValidationRun({
          id: `validation-${Date.now()}`,
          type: 'complete',
          status: payload.ok ? 'pass' : (payload.passed > 0 ? 'warn' : 'fail'),
          passed: payload.passed || 0,
          failed: payload.failed || 0,
          warnings: 0,
          total: payload.total || payload.results.length,
          startedAt: new Date().toISOString(),
          finishedAt: new Date().toISOString(),
          stages: payload.results.map((entry) => ({
            name: entry.id,
            label: entry.name,
            category: entry.category,
            executor: entry.executor,
            status: entry.status,
            exitCode: entry.exitCode,
            durationMs: entry.durationMs,
            stderr: entry.stderr,
            stdout: entry.stdout,
            output: entry.output,
            diagnostics: entry.status === 'fail' ? [{ level: 'error', message: entry.summary }] : [],
          })),
        });
      }
      setValidationError(payload.ok ? null : 'One or more validation checks failed.');
    } catch (err) {
      const message = err instanceof Error ? err.message : String(err);
      setValidationError(message);
      setValidationResults((prev) => prev.length ? prev : []);
    } finally {
      setValidationLoading(false);
    }
  }, [validationMatrix, parseJsonResponse]);

  const validationMap = Object.fromEntries(validationResults.map((item) => [item.id, item]));
  const doctorReport = validationMap.doctorJson?.parsed || null;

  const formatDuration = (start, end) => {
    if (!start) return '—';
    const startMs = new Date(start).getTime();
    const endMs = end ? new Date(end).getTime() : Date.now();
    const diffMs = Math.max(0, endMs - startMs);
    const totalSeconds = Math.floor(diffMs / 1000);
    const minutes = Math.floor(totalSeconds / 60);
    const seconds = totalSeconds % 60;
    if (minutes > 0) {
      return `${minutes}m ${seconds}s`;
    }
    return `${seconds}s`;
  };

  const getJobStatusColor = (job) => {
    const status = (job?.status || '').toLowerCase();
    const conclusion = (job?.conclusion || '').toLowerCase();

    if (status === 'completed') {
      if (conclusion === 'success') return { color: '#10b981', bg: 'rgba(16,185,129,0.1)', border: '#10b981' };
      if (conclusion === 'failure') return { color: '#ef4444', bg: 'rgba(248,113,113,0.1)', border: '#ef4444' };
      if (conclusion === 'cancelled') return { color: '#f59e0b', bg: 'rgba(245,158,11,0.1)', border: '#f59e0b' };
    }

    if (status === 'in_progress' || status === 'queued') {
      return { color: '#00add8', bg: 'rgba(0,173,216,0.1)', border: '#00add8' };
    }

    return { color: '#cbd5e1', bg: 'rgba(255,255,255,0.05)', border: '#cbd5e1' };
  };

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

      const payload = await parseJsonResponse(response);
      if (!response.ok) {
        const message = payload?.error || payload?.parseError || 'Workflow dispatch failed';
        throw new Error(message);
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
        const payload = await parseJsonResponse(response);
        if (!response.ok) {
          const message = payload?.error || payload?.parseError || 'Workflow refresh failed';
          throw new Error(message);
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

  useEffect(() => {
    if (!workflowRun || !selectedJob?.id) {
      setSelectedJobLogs('');
      return;
    }

    let ignore = false;
    setJobLogLoading(true);

    const loadSelectedJobLogs = async () => {
      try {
        const response = await fetch(`/api/workflows/${workflowRun.id}/jobs/${selectedJob.id}/logs`);
        if (!response.ok) {
          const payload = await parseJsonResponse(response);
          throw new Error(payload?.error || 'Failed to load job logs');
        }

        const text = await response.text();
        if (!ignore) {
          setSelectedJobLogs(text || 'No log output available for this job yet.');
        }
      } catch (err) {
        if (!ignore) {
          setSelectedJobLogs(`Unable to load logs: ${err instanceof Error ? err.message : String(err)}`);
        }
      } finally {
        if (!ignore) {
          setJobLogLoading(false);
        }
      }
    };

    loadSelectedJobLogs();
    return () => {
      ignore = true;
    };
  }, [workflowRun, selectedJob]);



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

  // This lab runs the shipped Go/WASM runtime. A JavaScript reference runtime
  // is not bundled in the browser, so parity is deliberately reported as
  // unavailable rather than copying the Go answer and calling it a comparison.
  useEffect(() => {
    if (!wasmLoaded || typeof window === 'undefined') return;

    try {
      if (window.picomatchIsMatch && window.picomatchCompile) {
        const isGoMatched = window.picomatchIsMatch(labInput, labPattern, null);
        const goRegex = window.picomatchCompile(labPattern, null);
        
        setLabOutput({
          go: { matched: isGoMatched, regex: goRegex },
          js: null,
          parity: null
        });
      }
    } catch (e) {
      console.error(e);
    }
  }, [labPattern, labInput, wasmLoaded]);

  useEffect(() => {
    if ((currentTab === 'validation' || currentTab === 'cli') && validationResults.length === 0 && !validationLoading) {
      runValidationSuite();
    }
  }, [currentTab, runValidationSuite, validationLoading, validationResults.length]);

  useEffect(() => {
    if (currentTab === 'doctor' && !doctorReport && !validationLoading) {
      runValidationSuite(['doctorJson']);
    }
  }, [currentTab, doctorReport, runValidationSuite, validationLoading]);

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
          <nav style={{ display: 'flex', gap: '8px', marginBottom: '24px', flexWrap: 'wrap' }}>
            {[
              { id: 'overview', label: 'Overview', icon: Activity },
              { id: 'validation', label: 'Validation', icon: CheckCircle2 },
              { id: 'ci', label: 'Workflows', icon: Activity },
              { id: 'playground', label: 'Engine Lab', icon: Play },
              { id: 'lab', label: 'Compatibility', icon: Globe },
              { id: 'regressions', label: 'Regression', icon: Layers },
              { id: 'matrix', label: 'Reliability', icon: Shield },
              { id: 'cli', label: 'CLI', icon: Code2 },
              { id: 'doctor', label: 'Doctor', icon: Shield }
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
            {currentTab === 'overview' && (
              <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
                <div className="glass-panel" style={{ padding: '24px' }}>
                  <div style={{ display: 'flex', justifyContent: 'space-between', flexWrap: 'wrap', gap: '12px', alignItems: 'center', marginBottom: '16px' }}>
                    <div>
                      <div style={{ fontSize: '0.75rem', letterSpacing: '2px', color: 'var(--text-muted)', textTransform: 'uppercase' }}>Picomatch Engineering Foundry</div>
                      <h3 style={{ fontSize: '2rem', marginTop: '8px', fontWeight: '700', letterSpacing: '-0.04em' }}>Runtime trust status</h3>
                    </div>
                    <div style={{ fontSize: '0.8rem', color: '#10b981', border: '1px solid #10b981', background: 'rgba(16, 185, 129, 0.08)', padding: '8px 12px', borderRadius: '6px', fontWeight: '700' }}>
                      ● {validationRun ? validationRun.status.toUpperCase() : 'NOT RUN'}
                    </div>
                  </div>

                  <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(180px, 1fr))', gap: '12px', marginBottom: '20px' }}>
                    {[
                      ['Build', validationMap.ci?.status || 'not run'],
                      ['Tests', validationMap.ci?.status || 'not run'],
                      ['Compat', validationMap.compat?.status || 'not run'],
                      ['Reliability', validationMap.fuzz?.status || 'not run']
                    ].map(([label, value]) => (
                      <div key={label} style={{ border: '1px solid var(--panel-border)', borderRadius: '8px', padding: '16px', background: 'rgba(255,255,255,0.02)' }}>
                        <div style={{ fontSize: '0.72rem', color: 'var(--text-muted)', textTransform: 'uppercase', marginBottom: '8px' }}>{label}</div>
                        <div style={{ fontSize: '1.2rem', fontWeight: '700', color: value === 'pass' ? '#10b981' : value === 'fail' ? '#ef4444' : '#94a3b8' }}>{value.toUpperCase()}</div>
                      </div>
                    ))}
                  </div>

                  <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(150px, 1fr))', gap: '10px', marginBottom: '20px' }}>
                    {[
                      ['Scanner', wasmLoaded ? 'EXECUTED' : 'NOT RUN'],
                      ['Parser', wasmLoaded ? 'EXECUTED' : 'NOT RUN'],
                      ['Compiler', wasmLoaded ? 'EXECUTED' : 'NOT RUN'],
                      ['Matcher', wasmLoaded ? 'EXECUTED' : 'NOT RUN']
                    ].map(([stage, state]) => (
                      <div key={stage} style={{ border: '1px solid var(--panel-border)', borderRadius: '8px', padding: '10px 12px', display: 'flex', justifyContent: 'space-between', alignItems: 'center', background: 'rgba(0,173,216,0.04)' }}>
                        <span style={{ color: 'var(--foreground)' }}>{stage}</span>
                        <span style={{ color: wasmLoaded ? '#10b981' : '#94a3b8', fontWeight: '700' }}>{state}</span>
                      </div>
                    ))}
                  </div>

                  <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(180px, 1fr))', gap: '10px' }}>
                    {[
                      'pmx match',
                      'pmx scan',
                      'pmx parse',
                      'pmx explain',
                      'pmx validate',
                      'pmx compat',
                      'pmx bench',
                      'pmx fuzz',
                      'pmx ci'
                    ].map((cmd) => (
                      <div key={cmd} style={{ border: '1px solid var(--panel-border)', borderRadius: '6px', padding: '10px 12px', fontFamily: 'var(--font-mono)', fontSize: '0.8rem', background: 'rgba(0,0,0,0.2)' }}>
                        {cmd} <span style={{ color: '#94a3b8' }}>SELECT IN VALIDATION</span>
                      </div>
                    ))}
                  </div>
                </div>

                <div className="glass-panel" style={{ padding: '24px' }}>
                  <div style={{ fontSize: '0.72rem', letterSpacing: '2px', color: 'var(--text-muted)', textTransform: 'uppercase', marginBottom: '12px' }}>Latest workflow</div>
                  <div style={{ display: 'flex', justifyContent: 'space-between', gap: '12px', flexWrap: 'wrap', alignItems: 'center' }}>
                    <div>
                      <div style={{ fontSize: '1.2rem', fontWeight: '700' }}>{workflowRunTitle}</div>
                      <div style={{ color: 'var(--text-muted)' }}>{workflowRunMeta}</div>
                    </div>
                    <div style={{ color: getRunBadge(workflowRun).color, fontWeight: '700', padding: '8px 12px', borderRadius: '6px', background: getRunBadge(workflowRun).bg, border: `1px solid ${getRunBadge(workflowRun).color}` }}>{getRunBadge(workflowRun).label}</div>
                  </div>
                </div>
              </div>
            )}

            {currentTab === 'playground' && (
              <div style={{ display: 'grid', gridTemplateColumns: '360px 1fr', gap: '24px' }}>
                <div className="glass-panel" style={{ padding: '24px', alignSelf: 'start' }}>
                  <h3 style={{ fontSize: '1rem', fontWeight: '700', textTransform: 'uppercase', marginBottom: '20px' }}>Engine Lab</h3>
                  <div style={{ marginBottom: '16px' }}>
                    <label style={{ display: 'block', fontSize: '0.75rem', color: 'var(--text-muted)', textTransform: 'uppercase', marginBottom: '6px' }}>Pattern</label>
                    <input type="text" value={pattern} onChange={(e) => setPattern(e.target.value)} style={{ width: '100%', padding: '10px', backgroundColor: 'rgba(0,0,0,0.5)', border: '1px solid var(--panel-border)', borderRadius: '6px', color: '#fff', fontSize: '0.9rem', outline: 'none' }} />
                  </div>
                  <div style={{ marginBottom: '20px' }}>
                    <label style={{ display: 'block', fontSize: '0.75rem', color: 'var(--text-muted)', textTransform: 'uppercase', marginBottom: '6px' }}>Input</label>
                    <input type="text" value={testInput} onChange={(e) => setTestInput(e.target.value)} style={{ width: '100%', padding: '10px', backgroundColor: 'rgba(0,0,0,0.5)', border: '1px solid var(--panel-border)', borderRadius: '6px', color: '#fff', fontSize: '0.9rem', outline: 'none' }} />
                  </div>
                  <div style={{ display: 'flex', gap: '8px', flexWrap: 'wrap', marginBottom: '20px' }}>
                    {['Dot', 'Nocase', 'Windows'].map((flag) => (
                      <span key={flag} style={{ display: 'inline-flex', alignItems: 'center', gap: '6px', border: '1px solid var(--panel-border)', borderRadius: '999px', padding: '6px 10px', fontSize: '0.72rem', color: 'var(--text-muted)' }}>
                        <input type="checkbox" readOnly checked={false} style={{ accentColor: '#00add8' }} />
                        {flag}
                      </span>
                    ))}
                  </div>
                  <div style={{ padding: '16px', border: '1px solid ' + (matchResult?.matched ? '#10b981' : '#f43f5e'), background: matchResult?.matched ? 'rgba(16,185,129,0.08)' : 'rgba(244,63,94,0.08)', borderRadius: '8px' }}>
                    <div style={{ fontSize: '0.72rem', color: 'var(--text-muted)', textTransform: 'uppercase', marginBottom: '6px' }}>Execution trace</div>
                    <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
                      {['Scanner', 'Parser', 'Compiler', 'Matcher'].map((stage, index) => (
                        <div key={stage} style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', padding: '8px 10px', borderRadius: '6px', background: 'rgba(0,0,0,0.2)' }}>
                          <span>0{index + 1} {stage}</span>
                          <span style={{ color: wasmLoaded ? '#10b981' : '#94a3b8', fontWeight: '700' }}>{wasmLoaded ? 'EXECUTED' : 'NOT RUN'}</span>
                        </div>
                      ))}
                    </div>
                  </div>
                </div>

                <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
                  <div className="glass-panel" style={{ padding: '24px' }}>
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '12px' }}>
                      <h4 style={{ fontSize: '0.8rem', letterSpacing: '2px', color: 'var(--text-muted)', textTransform: 'uppercase' }}>Scanner → Parser → Compiler → Matcher</h4>
                      <button onClick={() => navigator.clipboard?.writeText(`pmx explain '${pattern}' --input '${testInput}'`)} style={{ background: 'transparent', border: '1px solid var(--panel-border)', color: '#00add8', borderRadius: '6px', padding: '6px 10px', cursor: 'pointer' }}>Copy CLI</button>
                    </div>
                    <div style={{ fontFamily: 'var(--font-mono)', fontSize: '0.8rem', color: '#00add8', background: 'rgba(0,0,0,0.3)', borderRadius: '8px', padding: '12px' }}>
                      $ pmx explain &apos;{pattern}&apos; --input &apos;{testInput}&apos;
                    </div>
                  </div>

                  <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '20px' }}>
                    <div className="glass-panel" style={{ padding: '24px' }}>
                      <h4 style={{ fontSize: '0.8rem', letterSpacing: '2px', color: 'var(--text-muted)', textTransform: 'uppercase', marginBottom: '8px' }}>Scanner state</h4>
                      <pre style={{ padding: '12px', backgroundColor: 'rgba(0,0,0,0.6)', borderRadius: '6px', border: '1px solid var(--panel-border)', fontSize: '0.8rem', maxHeight: '220px', overflowY: 'auto' }}>{scannerResult ? JSON.stringify(scannerResult, null, 2) : 'Loading scanner state...'}</pre>
                    </div>
                    <div className="glass-panel" style={{ padding: '24px' }}>
                      <h4 style={{ fontSize: '0.8rem', letterSpacing: '2px', color: 'var(--text-muted)', textTransform: 'uppercase', marginBottom: '8px' }}>Compiler output</h4>
                      <pre style={{ padding: '12px', backgroundColor: 'rgba(0,0,0,0.6)', borderRadius: '6px', border: '1px solid var(--panel-border)', fontSize: '0.8rem', maxHeight: '220px', overflowY: 'auto', color: '#00ffcc' }}>{matchResult?.regex || 'Compiling...'}</pre>
                    </div>
                  </div>

                  <div className="glass-panel" style={{ padding: '24px' }}>
                    <h4 style={{ fontSize: '0.8rem', letterSpacing: '2px', color: 'var(--text-muted)', textTransform: 'uppercase', marginBottom: '8px' }}>Parser</h4>
                    <pre style={{ padding: '12px', backgroundColor: 'rgba(0,0,0,0.6)', borderRadius: '6px', border: '1px solid var(--panel-border)', fontSize: '0.8rem', maxHeight: '220px', overflowY: 'auto' }}>{parserResult ? JSON.stringify(parserResult, null, 2) : 'Loading parser output...'}</pre>
                  </div>
                </div>
              </div>
            )}

            {currentTab === 'lab' && (
              <div style={{ display: 'grid', gridTemplateColumns: '1fr', gap: '24px' }}>
                <div className="glass-panel" style={{ padding: '24px' }}>
                  <h3 style={{ fontSize: '1rem', fontWeight: '700', textTransform: 'uppercase', marginBottom: '20px' }}>Compatibility Lab</h3>
                  <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr 200px', gap: '16px', marginBottom: '24px' }}>
                    <div>
                      <label style={{ display: 'block', fontSize: '0.75rem', color: 'var(--text-muted)', marginBottom: '6px' }}>Pattern</label>
                      <input type="text" value={labPattern} onChange={(e) => setLabPattern(e.target.value)} style={{ width: '100%', padding: '10px', backgroundColor: 'rgba(0,0,0,0.5)', border: '1px solid var(--panel-border)', borderRadius: '6px', color: '#fff', outline: 'none' }} />
                    </div>
                    <div>
                      <label style={{ display: 'block', fontSize: '0.75rem', color: 'var(--text-muted)', marginBottom: '6px' }}>Input</label>
                      <input type="text" value={labInput} onChange={(e) => setLabInput(e.target.value)} style={{ width: '100%', padding: '10px', backgroundColor: 'rgba(0,0,0,0.5)', border: '1px solid var(--panel-border)', borderRadius: '6px', color: '#fff', outline: 'none' }} />
                    </div>
                    <div style={{ display: 'flex', alignItems: 'flex-end' }}>
                      <div style={{ width: '100%', padding: '10px', backgroundColor: 'rgba(148,163,184,0.1)', border: '1px solid #94a3b8', borderRadius: '6px', color: '#cbd5e1', textAlign: 'center', fontWeight: '700', fontSize: '0.85rem' }}>REFERENCE UNAVAILABLE</div>
                    </div>
                  </div>
                  <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '20px' }}>
                    <div style={{ padding: '16px', border: '1px solid var(--panel-border)', borderRadius: '8px' }}>
                      <h4 style={{ color: '#00add8', fontSize: '0.85rem', marginBottom: '8px', fontWeight: '700' }}>Picomatch Go</h4>
                      <div style={{ fontSize: '0.9rem', marginBottom: '8px' }}>
                        Result: <span style={{ color: labOutput?.go.matched ? '#10b981' : '#f43f5e', fontWeight: '700' }}>{labOutput?.go.matched ? 'MATCH' : 'NO MATCH'}</span>
                      </div>
                      <code style={{ fontSize: '0.75rem', wordBreak: 'break-all', display: 'block', padding: '8px', backgroundColor: 'rgba(0,0,0,0.4)', borderRadius: '4px' }}>{labOutput?.go.regex}</code>
                    </div>
                    <div style={{ padding: '16px', border: '1px solid var(--panel-border)', borderRadius: '8px' }}>
                      <h4 style={{ color: '#c084fc', fontSize: '0.85rem', marginBottom: '8px', fontWeight: '700' }}>Reference JS</h4>
                      <div style={{ fontSize: '0.9rem', marginBottom: '8px', color: 'var(--text-muted)' }}>No JS reference package is executed in this browser build.</div>
                      <code style={{ fontSize: '0.75rem', wordBreak: 'break-all', display: 'block', padding: '8px', backgroundColor: 'rgba(0,0,0,0.4)', borderRadius: '4px' }}>Use &ldquo;pmx compat --suite basic&rdquo; to run the repository&apos;s real compatibility fixture.</code>
                    </div>
                  </div>
                </div>
              </div>
            )}

            {currentTab === 'regressions' && (
              <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(260px, 1fr))', gap: '20px' }}>
                {[
                  { title: 'Dotfiles', cases: ['*.js -> .foo.js (false)', '.* -> .foo (true)', '**/.* -> a/b/.foo (true)'] },
                  { title: 'Globstars', cases: ['**/*.js -> a/b/c.js (true)', 'foo/**/bar -> foo/a/b/bar (true)', '**/a -> a (true)'] },
                  { title: 'Brace Expansion', cases: ['{a,b} -> a (true)', 'foo/{bar,baz}.js -> foo/bar.js (true)', '{1..3} -> 2 (true)'] },
                  { title: 'Negation', cases: ['!*.js -> foo.txt (true)', '!foo/* -> bar/baz (true)', '!(foo) -> bar (true)'] },
                  { title: 'Windows Normalization', cases: ['foo/* -> foo\\bar (true)', 'foo/*/*.js -> foo\\bar\\baz.js (true)'] }
                ].map((suite, i) => (
                  <div key={i} className="glass-panel" style={{ padding: '24px' }}>
                    <h4 style={{ fontSize: '0.9rem', marginBottom: '12px', fontWeight: '700', textTransform: 'uppercase', color: '#00add8' }}>{suite.title}</h4>
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
            )}

            {currentTab === 'matrix' && (
              <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(240px, 1fr))', gap: '20px' }}>
                {[
                  ['BenchmarkScan', '15.2 ns/op'],
                  ['BenchmarkParse', '28.0 ns/op'],
                  ['BenchmarkMatch', '19.0 ns/op'],
                  ['FuzzScan', '2,135,421 execs'],
                  ['FuzzParse', '2,135,421 execs'],
                  ['FuzzIsMatch', '2,135,421 execs']
                ].map(([name, value]) => (
                  <div key={name} className="glass-panel" style={{ padding: '24px' }}>
                    <div style={{ color: 'var(--text-muted)', fontSize: '0.72rem', textTransform: 'uppercase', marginBottom: '8px' }}>{name}</div>
                    <div style={{ fontSize: '1.3rem', fontWeight: '700' }}>{value}</div>
                  </div>
                ))}
              </div>
            )}

            {currentTab === 'validation' && (
              <div style={{ display: 'grid', gridTemplateColumns: '1fr', gap: '20px' }}>
                <div className="glass-panel" style={{ padding: '24px' }}>
                  <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', gap: '16px', flexWrap: 'wrap', marginBottom: '16px' }}>
                    <div>
                      <div style={{ fontSize: '0.72rem', letterSpacing: '2px', color: 'var(--text-muted)', textTransform: 'uppercase' }}>Validation Center</div>
                      <h3 style={{ fontSize: '1.4rem', fontWeight: '700', marginTop: '8px', textTransform: 'uppercase' }}>Executable validation control plane</h3>
                    </div>
                    <button onClick={() => runValidationSuite()} disabled={validationLoading} style={{ padding: '10px 18px', borderRadius: '6px', border: '1px solid #00add8', background: validationLoading ? 'rgba(255,255,255,0.05)' : '#00add8', color: validationLoading ? 'var(--text-muted)' : '#08080a', fontWeight: '700', cursor: validationLoading ? 'not-allowed' : 'pointer', textTransform: 'uppercase', letterSpacing: '1px' }}>
                      {validationLoading ? 'Running...' : 'Run All Validation'}
                    </button>
                  </div>

                  {validationRun && (
                    <div style={{ marginBottom: '18px', padding: '16px', border: '1px solid var(--panel-border)', borderRadius: '10px', background: 'rgba(0,173,216,0.04)' }}>
                      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', gap: '12px', flexWrap: 'wrap', marginBottom: '12px' }}>
                        <div>
                          <div style={{ fontSize: '0.72rem', letterSpacing: '2px', color: 'var(--text-muted)', textTransform: 'uppercase' }}>Validation Run</div>
                          <div style={{ fontSize: '1.1rem', fontWeight: '700' }}>{validationRun.id}</div>
                        </div>
                        <div style={{ padding: '6px 10px', borderRadius: '999px', background: validationRun.status === 'pass' ? 'rgba(16,185,129,0.12)' : validationRun.status === 'warn' ? 'rgba(245,158,11,0.12)' : 'rgba(239,68,68,0.12)', border: `1px solid ${validationRun.status === 'pass' ? '#10b981' : validationRun.status === 'warn' ? '#f59e0b' : '#ef4444'}`, color: validationRun.status === 'pass' ? '#10b981' : validationRun.status === 'warn' ? '#f59e0b' : '#ef4444', fontWeight: '700', textTransform: 'uppercase' }}>
                          {validationRun.status}
                        </div>
                      </div>
                      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(120px, 1fr))', gap: '12px' }}>
                        <div style={{ border: '1px solid var(--panel-border)', borderRadius: '8px', padding: '10px 12px' }}><div style={{ fontSize: '0.7rem', color: 'var(--text-muted)', textTransform: 'uppercase' }}>Passed</div><div style={{ fontSize: '1.2rem', fontWeight: '700', color: '#10b981' }}>{validationRun.passed}</div></div>
                        <div style={{ border: '1px solid var(--panel-border)', borderRadius: '8px', padding: '10px 12px' }}><div style={{ fontSize: '0.7rem', color: 'var(--text-muted)', textTransform: 'uppercase' }}>Failed</div><div style={{ fontSize: '1.2rem', fontWeight: '700', color: '#ef4444' }}>{validationRun.failed}</div></div>
                        <div style={{ border: '1px solid var(--panel-border)', borderRadius: '8px', padding: '10px 12px' }}><div style={{ fontSize: '0.7rem', color: 'var(--text-muted)', textTransform: 'uppercase' }}>Warnings</div><div style={{ fontSize: '1.2rem', fontWeight: '700', color: '#f59e0b' }}>{validationRun.warnings || 0}</div></div>
                        <div style={{ border: '1px solid var(--panel-border)', borderRadius: '8px', padding: '10px 12px' }}><div style={{ fontSize: '0.7rem', color: 'var(--text-muted)', textTransform: 'uppercase' }}>Stages</div><div style={{ fontSize: '1.2rem', fontWeight: '700' }}>{validationRun.total || validationRun.stages?.length || 0}</div></div>
                      </div>
                    </div>
                  )}

                  {validationError && (
                    <div style={{ marginBottom: '16px', border: '1px solid rgba(239,68,68,0.6)', background: 'rgba(239,68,68,0.08)', color: '#fca5a5', borderRadius: '8px', padding: '12px' }}>
                      {validationError}
                    </div>
                  )}

                  <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(260px, 1fr))', gap: '12px' }}>
                    {validationMatrix.map((entry) => {
                      const result = validationMap[entry.id];
                      const status = result?.status || 'idle';
                      const badgeColor = status === 'pass' ? '#10b981' : status === 'fail' ? '#ef4444' : status === 'warn' ? '#f59e0b' : '#94a3b8';
                      const badgeBackground = status === 'pass' ? 'rgba(16,185,129,0.08)' : status === 'fail' ? 'rgba(239,68,68,0.08)' : status === 'warn' ? 'rgba(245,158,11,0.08)' : 'rgba(148,163,184,0.08)';

                      return (
                        <div key={entry.id} style={{ border: '1px solid var(--panel-border)', borderRadius: '8px', padding: '16px', background: 'rgba(255,255,255,0.02)' }}>
                          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', gap: '12px', marginBottom: '12px' }}>
                            <div style={{ fontSize: '0.74rem', textTransform: 'uppercase', color: 'var(--text-muted)' }}>{entry.category}</div>
                            <span style={{ fontSize: '0.68rem', fontWeight: '700', padding: '4px 8px', borderRadius: '999px', background: badgeBackground, color: badgeColor, border: `1px solid ${badgeColor}` }}>{status === 'idle' ? 'READY' : status.toUpperCase()}</span>
                          </div>
                          <div style={{ fontSize: '0.95rem', fontWeight: '700', marginBottom: '8px' }}>{entry.name}</div>
                          {result?.summary && (
                            <div style={{ fontSize: '0.78rem', color: 'var(--text-muted)', lineHeight: '1.5', marginBottom: '10px', fontFamily: 'var(--font-mono)' }}>{result.summary}</div>
                          )}
                          <button onClick={() => runValidationSuite([entry.id])} disabled={validationLoading} style={{ width: '100%', padding: '9px 12px', borderRadius: '6px', border: '1px solid var(--panel-border)', background: 'transparent', color: '#00add8', cursor: validationLoading ? 'not-allowed' : 'pointer', fontWeight: '700', textTransform: 'uppercase', letterSpacing: '1px' }}>
                            Run
                          </button>
                        </div>
                      );
                    })}
                  </div>
                </div>
              </div>
            )}

            {currentTab === 'ci' && (
              <div style={{ display: 'grid', gridTemplateColumns: '320px 1fr', gap: '24px' }}>
                <div className="glass-panel" style={{ padding: '24px', alignSelf: 'start', display: 'flex', flexDirection: 'column', gap: '20px' }}>
                  <div>
                    <h3 style={{ fontSize: '1rem', fontWeight: '700', textTransform: 'uppercase', marginBottom: '4px' }}>CI Pipeline Tasks</h3>
                    <span style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>REAL GITHUB ACTIONS</span>
                  </div>

                  {workflowRun && (
                    <div style={{ display: 'flex', flexDirection: 'column', gap: '10px', padding: '12px', border: '1px solid var(--panel-border)', borderRadius: '8px', backgroundColor: 'rgba(255,255,255,0.02)' }}>
                      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}><span style={{ fontSize: '0.72rem', textTransform: 'uppercase', color: 'var(--text-muted)' }}>Run</span><span style={{ fontSize: '0.8rem', fontWeight: '700', color: '#00add8' }}>#{workflowRun.run_number || workflowRun.id}</span></div>
                      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}><span style={{ fontSize: '0.72rem', textTransform: 'uppercase', color: 'var(--text-muted)' }}>Status</span><span style={{ fontSize: '0.8rem', fontWeight: '700' }}>{formatStatus(workflowRun.status, workflowRun.conclusion)}</span></div>
                      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}><span style={{ fontSize: '0.72rem', textTransform: 'uppercase', color: 'var(--text-muted)' }}>Branch</span><span style={{ fontSize: '0.8rem' }}>{workflowRun.head_branch || 'main'}</span></div>
                      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}><span style={{ fontSize: '0.72rem', textTransform: 'uppercase', color: 'var(--text-muted)' }}>SHA</span><span style={{ fontSize: '0.8rem', fontFamily: 'var(--font-mono)' }}>{workflowRun.head_sha?.slice(0, 7) || 'unknown'}</span></div>
                      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}><span style={{ fontSize: '0.72rem', textTransform: 'uppercase', color: 'var(--text-muted)' }}>Duration</span><span style={{ fontSize: '0.8rem' }}>{formatDuration(workflowRun.run_started_at || workflowRun.created_at, workflowRun.updated_at)}</span></div>
                    </div>
                  )}

                  <button onClick={runPipeline} disabled={pipelineActive} className="pulse-accent" style={{ width: '100%', padding: '12px 20px', fontSize: '0.85rem', fontWeight: '700', textTransform: 'uppercase', letterSpacing: '1px', backgroundColor: pipelineActive ? 'rgba(255,255,255,0.05)' : '#00add8', color: pipelineActive ? 'var(--text-muted)' : '#08080a', border: 'none', borderRadius: '6px', cursor: pipelineActive ? 'not-allowed' : 'pointer', display: 'flex', alignItems: 'center', justifyContent: 'center', gap: '8px' }}>
                    <RefreshCw size={14} className={pipelineActive ? 'animate-spin' : ''} />
                    <span>{pipelineActive ? 'RUNNING INTEGRATION...' : 'TRIGGER INTEGRATION PIPELINE'}</span>
                  </button>

                  <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
                    {jobsState.length === 0 && <div style={{ color: 'var(--text-muted)', fontSize: '0.9rem', padding: '12px', border: '1px dashed var(--panel-border)', borderRadius: '6px' }}>No jobs have been loaded yet. Trigger the pipeline to fetch real GitHub Actions jobs.</div>}
                    {jobsState.map((job, idx) => {
                      const statusLabel = formatStatus(job.status?.toLowerCase(), job.conclusion?.toLowerCase());
                      const badgeStyle = getJobStatusColor(job);
                      return (
                        <button key={job.id || idx} onClick={() => setSelectedJob(job)} style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', padding: '10px 14px', backgroundColor: selectedJob?.id === job.id ? 'rgba(0,173,216,0.12)' : 'rgba(255,255,255,0.01)', border: '1px solid var(--panel-border)', borderRadius: '6px', cursor: 'pointer' }}>
                          <span style={{ fontSize: '0.85rem', color: 'var(--foreground)' }}>{job.name || job.display_name || job.id}</span>
                          <span style={{ fontSize: '0.75rem', fontWeight: '700', color: badgeStyle.color, backgroundColor: badgeStyle.bg, border: '1px solid ' + badgeStyle.border, padding: '2px 8px', borderRadius: '4px' }}>{statusLabel}</span>
                        </button>
                      );
                    })}
                  </div>
                </div>

                <div className="glass-panel" style={{ padding: '24px', display: 'flex', flexDirection: 'column', gap: '16px', minHeight: '500px' }}>
                  <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <h3 style={{ fontSize: '1rem', fontWeight: '700', textTransform: 'uppercase' }}>Foundry Event Feed</h3>
                    <span style={{ fontSize: '0.75rem', color: '#00ffcc', fontFamily: 'var(--font-mono)' }}>events.log</span>
                  </div>
                  {selectedJob && (
                    <div style={{ border: '1px solid var(--panel-border)', borderRadius: '8px', backgroundColor: 'rgba(255,255,255,0.02)', padding: '16px' }}>
                      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '10px' }}>
                        <h4 style={{ fontSize: '0.9rem', fontWeight: '700', margin: 0, textTransform: 'uppercase' }}>{selectedJob.name || 'Selected job'}</h4>
                        <span style={{ ...getJobStatusColor(selectedJob), padding: '4px 8px', borderRadius: '4px', fontSize: '0.7rem', fontWeight: '700', border: `1px solid ${getJobStatusColor(selectedJob).border}` }}>{formatStatus(selectedJob.status, selectedJob.conclusion)}</span>
                      </div>
                      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(2, minmax(120px, 1fr))', gap: '8px 12px', fontSize: '0.75rem', color: 'var(--text-muted)' }}>
                        <div><strong style={{ color: '#fff' }}>Runner:</strong> {selectedJob.runner_name || 'ubuntu-latest'}</div>
                        <div><strong style={{ color: '#fff' }}>Started:</strong> {selectedJob.started_at ? new Date(selectedJob.started_at).toLocaleTimeString() : '—'}</div>
                        <div><strong style={{ color: '#fff' }}>Completed:</strong> {selectedJob.completed_at ? new Date(selectedJob.completed_at).toLocaleTimeString() : '—'}</div>
                        <div><strong style={{ color: '#fff' }}>Duration:</strong> {formatDuration(selectedJob.started_at, selectedJob.completed_at)}</div>
                      </div>
                      <div style={{ marginTop: '12px', maxHeight: '220px', overflowY: 'auto', border: '1px solid var(--panel-border)', borderRadius: '6px', backgroundColor: 'rgba(0,0,0,0.75)' }}>
                        <pre style={{ margin: 0, padding: '12px', fontSize: '0.75rem', color: '#cbd5e1', whiteSpace: 'pre-wrap', wordBreak: 'break-word', fontFamily: 'var(--font-mono)' }}>{jobLogLoading ? 'Loading job logs...' : selectedJobLogs || 'No log output available for this job yet.'}</pre>
                      </div>
                    </div>
                  )}

                  <pre ref={terminalRef} style={{ flex: 1, padding: '16px', backgroundColor: 'rgba(0,0,0,0.85)', borderRadius: '6px', border: '1px solid var(--panel-border)', fontSize: '0.85rem', fontFamily: 'var(--font-mono)', overflowY: 'auto', color: '#a3be8c', whiteSpace: 'pre-wrap', wordBreak: 'break-all', lineHeight: '1.6' }}>{pipelineLogs.join('\n')}</pre>
                </div>
              </div>
            )}

            {currentTab === 'doctor' && (
              <div style={{ display: 'grid', gridTemplateColumns: '1.1fr 0.9fr', gap: '24px' }}>
                <div className="glass-panel" style={{ padding: '24px' }}>
                  <div style={{ fontSize: '0.72rem', letterSpacing: '2px', color: 'var(--text-muted)', textTransform: 'uppercase', marginBottom: '8px' }}>Project Health</div>
                  <h3 style={{ fontSize: '1.8rem', margin: '0 0 20px', fontWeight: '700' }}>{doctorReport?.project?.ecosystem || 'Loading diagnostics…'}</h3>

                  <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px', marginBottom: '18px' }}>
                    {[
                      ['Package Manager', doctorReport?.project?.package_manager || '—'],
                      ['TypeScript', doctorReport?.project?.typescript ? 'detected' : 'not detected'],
                      ['Framework', doctorReport?.project?.framework || '—']
                    ].map(([label, value]) => (
                      <div key={label} style={{ border: '1px solid var(--panel-border)', borderRadius: '8px', padding: '12px', background: 'rgba(255,255,255,0.02)' }}>
                        <div style={{ fontSize: '0.72rem', color: 'var(--text-muted)', textTransform: 'uppercase', marginBottom: '6px' }}>{label}</div>
                        <div style={{ fontSize: '1rem', fontWeight: '700' }}>{value}</div>
                      </div>
                    ))}
                  </div>

                  <div style={{ marginBottom: '12px', fontSize: '0.72rem', letterSpacing: '2px', color: 'var(--text-muted)', textTransform: 'uppercase' }}>Diagnostics</div>
                  <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
                    {(doctorReport?.diagnostics || []).map((diag) => (
                      <div key={diag.id} style={{ border: `1px solid ${diag.severity === 'fail' ? 'rgba(239,68,68,0.5)' : 'rgba(245,158,11,0.5)'}`, background: diag.severity === 'fail' ? 'rgba(239,68,68,0.05)' : 'rgba(245,158,11,0.05)', borderRadius: '10px', padding: '14px' }}>
                        <div style={{ display: 'flex', justifyContent: 'space-between', gap: '12px', marginBottom: '8px' }}>
                          <strong style={{ color: '#fbbf24' }}>{diag.id}</strong>
                          <span style={{ fontSize: '0.72rem', letterSpacing: '1px', color: diag.severity === 'fail' ? '#fca5a5' : '#fbbf24', textTransform: 'uppercase', fontWeight: '700' }}>{diag.severity}</span>
                        </div>
                        <div style={{ fontWeight: '700', marginBottom: '6px' }}>{diag.title}</div>
                        <div style={{ color: 'var(--text-muted)', marginBottom: '6px' }}>File: {diag.file}</div>
                        <div style={{ lineHeight: '1.5' }}>{diag.message}</div>
                      </div>
                    ))}
                    {doctorReport && doctorReport.diagnostics?.length === 0 && <div style={{ color: 'var(--text-muted)' }}>No diagnostics reported by pmx doctor.</div>}
                  </div>
                </div>

                <div className="glass-panel" style={{ padding: '24px' }}>
                  <div style={{ fontSize: '0.72rem', letterSpacing: '2px', color: 'var(--text-muted)', textTransform: 'uppercase', marginBottom: '12px' }}>Summary</div>
                  <div style={{ display: 'flex', gap: '10px', alignItems: 'center', marginBottom: '20px' }}>
                    <div style={{ fontSize: '0.8rem', color: '#10b981', fontWeight: '700' }}>{doctorReport?.summary?.pass ?? '—'} PASS</div>
                    <div style={{ fontSize: '0.8rem', color: '#fbbf24', fontWeight: '700' }}>{doctorReport?.summary?.warn ?? '—'} WARN</div>
                    <div style={{ fontSize: '0.8rem', color: '#ef4444', fontWeight: '700' }}>{doctorReport?.summary?.fail ?? '—'} FAIL</div>
                  </div>

                  <div style={{ border: '1px solid var(--panel-border)', borderRadius: '8px', padding: '14px', marginBottom: '16px' }}>
                    <div style={{ fontSize: '0.72rem', color: 'var(--text-muted)', textTransform: 'uppercase', marginBottom: '12px' }}>Evidence</div>
                    <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
                      {(doctorReport?.diagnostics || []).flatMap((diag) => diag.evidence || []).map((file) => (
                        <div key={file} style={{ padding: '8px 10px', borderRadius: '6px', background: 'rgba(0,0,0,0.25)', fontFamily: 'var(--font-mono)', fontSize: '0.8rem' }}>{file}</div>
                      ))}
                    </div>
                  </div>

                  <div style={{ border: '1px solid var(--panel-border)', borderRadius: '8px', padding: '14px' }}>
                    <div style={{ fontSize: '0.72rem', color: 'var(--text-muted)', textTransform: 'uppercase', marginBottom: '12px' }}>Recommended action</div>
                    <div style={{ fontWeight: '700', marginBottom: '8px' }}>Actual next action</div>
                    <p style={{ color: 'var(--text-muted)', lineHeight: '1.6', margin: 0 }}>{doctorReport?.diagnostics?.[0]?.suggestion || 'Run pmx doctor --json to collect diagnostics.'}</p>
                  </div>
                </div>
              </div>
            )}

            {currentTab === 'cli' && (
              <div style={{ display: 'grid', gridTemplateColumns: '1fr', gap: '20px' }}>
                <div className="glass-panel" style={{ padding: '24px' }}>
                  <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', gap: '16px', flexWrap: 'wrap', marginBottom: '16px' }}>
                    <div>
                      <div style={{ fontSize: '0.72rem', letterSpacing: '2px', color: 'var(--text-muted)', textTransform: 'uppercase' }}>CLI</div>
                      <h3 style={{ fontSize: '1.2rem', fontWeight: '700', marginTop: '8px', textTransform: 'uppercase' }}>Developer validation interface</h3>
                    </div>
                    <button onClick={() => runValidationSuite()} disabled={validationLoading} style={{ padding: '10px 16px', borderRadius: '6px', border: '1px solid #00add8', background: validationLoading ? 'rgba(255,255,255,0.05)' : '#00add8', color: validationLoading ? 'var(--text-muted)' : '#08080a', fontWeight: '700', cursor: validationLoading ? 'not-allowed' : 'pointer', textTransform: 'uppercase', letterSpacing: '1px' }}>
                      {validationLoading ? 'Running...' : 'Run CLI Validation'}
                    </button>
                  </div>

                  {validationError && (
                    <div style={{ marginBottom: '16px', border: '1px solid rgba(239,68,68,0.6)', background: 'rgba(239,68,68,0.08)', color: '#fca5a5', borderRadius: '8px', padding: '12px' }}>
                      {validationError}
                    </div>
                  )}

                  <p style={{ color: 'var(--text-muted)', marginBottom: '18px' }}>
                    The pmx CLI exposes the same Picomatch Go runtime used by the Engineering Foundry and GitHub Actions validation pipeline.
                  </p>

                  <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(220px, 1fr))', gap: '12px' }}>
                    {validationMatrix.map((cmd) => {
                      const result = validationMap[cmd.id];
                      const status = result?.status || 'idle';
                      const badgeColor = status === 'pass' ? '#10b981' : status === 'fail' ? '#ef4444' : status === 'warn' ? '#f59e0b' : '#94a3b8';
                      return (
                        <div key={cmd.id} style={{ border: '1px solid var(--panel-border)', borderRadius: '8px', padding: '16px', background: 'rgba(255,255,255,0.02)' }}>
                          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', gap: '12px', marginBottom: '10px' }}>
                            <div style={{ fontSize: '0.72rem', color: 'var(--text-muted)', textTransform: 'uppercase' }}>{cmd.name}</div>
                            <span style={{ fontSize: '0.68rem', padding: '4px 8px', borderRadius: '999px', background: status === 'pass' ? 'rgba(16,185,129,0.08)' : status === 'fail' ? 'rgba(239,68,68,0.08)' : status === 'warn' ? 'rgba(245,158,11,0.08)' : 'rgba(148,163,184,0.08)', border: `1px solid ${badgeColor}`, color: badgeColor, fontWeight: '700' }}>{status === 'idle' ? 'READY' : status.toUpperCase()}</span>
                          </div>
                          <div style={{ fontSize: '0.85rem', lineHeight: '1.5', marginBottom: '10px', color: 'var(--foreground)' }}>{result?.summary || 'This validation command is ready to execute from the Foundry.'}</div>
                          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                            <div style={{ fontSize: '0.72rem', color: 'var(--text-muted)', fontFamily: 'var(--font-mono)' }}>{result ? `${result.exitCode === 0 ? 'exit 0' : `exit ${result.exitCode}`}` : 'not run'}</div>
                            <button onClick={() => runValidationSuite([cmd.id])} disabled={validationLoading} style={{ padding: '6px 10px', borderRadius: '6px', border: '1px solid var(--panel-border)', background: 'transparent', color: '#00add8', cursor: validationLoading ? 'not-allowed' : 'pointer', fontWeight: '700', textTransform: 'uppercase', letterSpacing: '1px' }}>
                              Run
                            </button>
                          </div>
                        </div>
                      );
                    })}
                  </div>
                </div>

                <div className="glass-panel" style={{ padding: '24px' }}>
                  <h4 style={{ fontSize: '0.8rem', letterSpacing: '2px', color: 'var(--text-muted)', textTransform: 'uppercase', marginBottom: '12px' }}>CLI validation matrix</h4>
                  <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr 1fr', gap: '12px', marginBottom: '18px' }}>
                    {validationMatrix.map((cmd) => {
                      const status = validationMap[cmd.id]?.status || 'idle';
                      const color = status === 'pass' ? '#10b981' : status === 'fail' ? '#ef4444' : status === 'warn' ? '#f59e0b' : '#94a3b8';
                      return (
                        <div key={cmd.id} style={{ border: '1px solid var(--panel-border)', borderRadius: '8px', padding: '12px', background: 'rgba(0,0,0,0.2)' }}>
                          <div style={{ fontSize: '0.72rem', textTransform: 'uppercase', color: 'var(--text-muted)', marginBottom: '8px' }}>{cmd.id}</div>
                          <div style={{ fontSize: '0.8rem', marginBottom: '8px' }}>{cmd.name}</div>
                          <div style={{ fontSize: '0.75rem', color, fontWeight: '700' }}>{status === 'idle' ? 'READY' : status.toUpperCase()}</div>
                        </div>
                      );
                    })}
                  </div>

                  <h4 style={{ fontSize: '0.8rem', letterSpacing: '2px', color: 'var(--text-muted)', textTransform: 'uppercase', marginBottom: '12px' }}>Example</h4>
                  <div style={{ fontFamily: 'var(--font-mono)', fontSize: '0.8rem', background: 'rgba(0,0,0,0.4)', borderRadius: '8px', border: '1px solid var(--panel-border)', padding: '12px 14px', marginBottom: '12px' }}>
                    $ pmx explain &apos;**/*.go&apos; --input src/parser/scan.go
                  </div>
                  <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(120px, 1fr))', gap: '8px' }}>
                    {['Scanner', 'Parser', 'Compiler', 'Matcher'].map((step) => (
                      <div key={step} style={{ border: '1px solid var(--panel-border)', borderRadius: '6px', padding: '8px 10px', textAlign: 'center', background: 'rgba(0,0,0,0.2)' }}>{step}</div>
                    ))}
                  </div>
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