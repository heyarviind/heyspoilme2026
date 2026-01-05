<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { api } from '$lib/api';
	import { auth } from '$lib/stores/auth';
	import Header from '$lib/components/Header.svelte';
	import Footer from '$lib/components/Footer.svelte';

	let step = $state(1);
	let loading = $state(true);
	let submitting = $state(false);
	let error = $state<string | null>(null);
	let verificationCode = $state<string | null>(null);
	let status = $state<string | null>(null);

	// Form data
	let documentType = $state<string>('aadhar');
	let documentFile = $state<File | null>(null);
	let documentPreview = $state<string | null>(null);
	let documentUrl = $state<string | null>(null);
	let videoFile = $state<File | null>(null);
	let videoPreview = $state<string | null>(null);
	let videoUrl = $state<string | null>(null);
	let uploadingDocument = $state(false);
	let uploadingVideo = $state(false);

	let authState = $state<any>(null);
	auth.subscribe(s => authState = s);

	const documentTypes = [
		{ value: 'aadhar', label: 'Aadhaar Card', icon: '🪪' },
		{ value: 'passport', label: 'Passport', icon: '📘' },
		{ value: 'driving_license', label: 'Driving License', icon: '🚗' },
	];

	const benefits = [
		{ icon: '✓', title: '3x More Matches', description: 'Verified profiles appear at the top and get significantly more attention' },
		{ icon: '🛡️', title: 'Build Trust', description: 'Show others you are who you say you are' },
		{ icon: '⭐', title: 'Verified Badge', description: 'Get a green checkmark on your profile' },
		{ icon: '🔒', title: 'Priority Support', description: 'Get faster responses from our support team' },
	];

	onMount(async () => {
		// Check if already verified
		if (authState?.profile?.is_verified) {
			goto('/profile');
			return;
		}

		// Check existing verification status
		try {
			const statusData = await api.getVerificationStatus() as { status: string; rejection_reason?: string };
			status = statusData.status;
			
			if (status === 'pending') {
				step = 5; // Show pending status
			} else if (status === 'approved') {
				goto('/profile');
				return;
			}
		} catch (e) {
			// No verification request yet, continue
		}

		// Generate verification code
		try {
			const codeData = await api.getVerificationCode() as { code: string };
			verificationCode = codeData.code;
		} catch (e) {
			error = 'Failed to generate verification code';
		}

		loading = false;
	});

	async function uploadFile(file: File, type: 'document' | 'video'): Promise<string> {
		const ext = file.name.split('.').pop()?.toLowerCase() || '';
		const contentType = file.type;

		console.log(`Uploading ${type}:`, { ext, contentType, size: file.size });

		const presignedData = await api.getPresignedUrl(ext, contentType) as { public_url: string; upload_url: string };
		console.log('Got presigned URL:', presignedData);

		try {
			const uploadResponse = await fetch(presignedData.upload_url, {
				method: 'PUT',
				body: file,
				headers: { 'Content-Type': contentType },
			});

			if (!uploadResponse.ok) {
				const errorText = await uploadResponse.text().catch(() => '');
				console.error('S3 upload failed:', uploadResponse.status, uploadResponse.statusText, errorText);
				throw new Error(`Upload failed: ${uploadResponse.status} ${uploadResponse.statusText}`);
			}
		} catch (fetchError: any) {
			console.error('S3 fetch error:', fetchError);
			// This might be a CORS error - try to detect it
			if (fetchError.message?.includes('Failed to fetch') || fetchError.name === 'TypeError') {
				throw new Error('Upload failed - CORS or network error. Check S3 bucket CORS settings.');
			}
			throw fetchError;
		}

		console.log('Upload successful, URL:', presignedData.public_url);
		return presignedData.public_url;
	}

	async function handleDocumentChange(e: Event) {
		const target = e.target as HTMLInputElement;
		const file = target.files?.[0];
		if (!file) return;

		documentFile = file;
		documentPreview = URL.createObjectURL(file);

		// Upload immediately
		uploadingDocument = true;
		error = null;
		try {
			documentUrl = await uploadFile(file, 'document');
			console.log('Document uploaded successfully:', documentUrl);
		} catch (e: any) {
			console.error('Document upload error:', e);
			error = e.message || 'Failed to upload document. Please try again.';
			documentFile = null;
			documentPreview = null;
			documentUrl = null;
		} finally {
			uploadingDocument = false;
		}
	}

	async function handleVideoChange(e: Event) {
		const target = e.target as HTMLInputElement;
		const file = target.files?.[0];
		if (!file) return;

		// Validate video size (max 50MB)
		if (file.size > 50 * 1024 * 1024) {
			error = 'Video must be less than 50MB';
			return;
		}

		videoFile = file;
		videoPreview = URL.createObjectURL(file);

		// Upload immediately
		uploadingVideo = true;
		error = null;
		try {
			videoUrl = await uploadFile(file, 'video');
		} catch (e: any) {
			error = 'Failed to upload video. Please try again.';
			videoFile = null;
			videoPreview = null;
		} finally {
			uploadingVideo = false;
		}
	}

	async function submitVerification() {
		if (!documentUrl || !videoUrl || !verificationCode) {
			error = 'Please complete all steps';
			return;
		}

		submitting = true;
		error = null;

		try {
			await api.submitVerification(documentType, documentUrl, videoUrl, verificationCode);
			step = 5; // Success
			status = 'pending';
		} catch (e: any) {
			error = e.message || 'Failed to submit verification';
		} finally {
			submitting = false;
		}
	}

	function nextStep() {
		if (step === 2 && !documentUrl) {
			error = 'Please upload a document';
			return;
		}
		if (step === 3 && !videoUrl) {
			error = 'Please record and upload your video';
			return;
		}
		error = null;
		step++;
	}

	function prevStep() {
		error = null;
		step--;
	}
</script>

<svelte:head>
	<title>Verify Your Identity | HeySpoilMe</title>
</svelte:head>

<div class="verify-page">
	<div class="desktop-header">
		<Header />
	</div>
	<a href="/browse" class="mobile-close" aria-label="Close">
		<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
			<path d="M18 6L6 18M6 6l12 12"/>
		</svg>
	</a>

	<main class="container">
		{#if loading}
			<div class="loading">
				<div class="spinner"></div>
				<p>Loading...</p>
			</div>
		{:else if step === 5}
			<!-- Pending/Success State -->
			<div class="status-card">
				<div class="status-icon pending">⏳</div>
				<h1>Verification Submitted</h1>
				<p>Thank you for submitting your verification request. Our team will review it within 24-48 hours.</p>
				<p class="hint">You'll receive a notification once your verification is complete.</p>
				<a href="/browse" class="btn-primary">Continue Browsing</a>
			</div>
		{:else}
			<div class="verify-content">
				<!-- Progress Bar -->
				<div class="progress-bar">
					{#each [1, 2, 3, 4] as s}
						<div class="step" class:active={step >= s} class:current={step === s}>
							<span class="step-number">{s}</span>
						</div>
						{#if s < 4}
							<div class="step-line" class:active={step > s}></div>
						{/if}
					{/each}
				</div>

				{#if step === 1}
					<!-- Step 1: Benefits -->
					<div class="step-content">
						<div class="header-section">
							<div class="shield-icon">🛡️</div>
							<h1>Get Verified</h1>
							<p class="subtitle">Unlock premium features and build trust with potential matches</p>
						</div>

						<div class="benefits-grid">
							{#each benefits as benefit}
								<div class="benefit-card">
									<span class="benefit-icon">{benefit.icon}</span>
									<h3>{benefit.title}</h3>
									<p>{benefit.description}</p>
								</div>
							{/each}
						</div>

					<div class="privacy-notice">
						<span class="lock-icon">🔐</span>
						<div>
							<strong>Your privacy is protected</strong>
							<p>We will <strong>permanently delete</strong> your documents and video immediately after verification. We do not store this information.</p>
						</div>
					</div>

					<div class="btn-group single">
						<button class="btn-primary" onclick={() => step = 2}>
							Start Verification
						</button>
					</div>
				</div>

				{:else if step === 2}
					<!-- Step 2: Document Upload -->
					<div class="step-content">
						<h2>Upload Identity Document</h2>
						<p class="step-description">Upload a clear photo of one of the following documents:</p>

						<div class="document-types">
							{#each documentTypes as doc}
								<label class="document-option" class:selected={documentType === doc.value}>
									<input type="radio" name="documentType" value={doc.value} bind:group={documentType} />
									<span class="doc-icon">{doc.icon}</span>
									<span class="doc-label">{doc.label}</span>
								</label>
							{/each}
						</div>

						<div class="upload-area">
							{#if documentPreview}
								<div class="preview-container">
									<img src={documentPreview} alt="Document preview" class="document-preview" />
									{#if uploadingDocument}
										<div class="upload-overlay">
											<div class="spinner"></div>
											<span>Uploading...</span>
										</div>
									{:else if documentUrl}
										<div class="upload-success">✓ Uploaded</div>
										<button class="remove-btn" onclick={() => { documentFile = null; documentPreview = null; documentUrl = null; }}>
											Change
										</button>
									{:else}
										<div class="upload-failed">✗ Upload failed</div>
										<button class="remove-btn" onclick={() => { documentFile = null; documentPreview = null; documentUrl = null; }}>
											Try Again
										</button>
									{/if}
								</div>
							{:else}
								<label class="upload-box">
									<input type="file" accept="image/*" onchange={handleDocumentChange} />
									<span class="upload-icon">📷</span>
									<span class="upload-text">Click to upload document photo</span>
									<span class="upload-hint">Make sure all details are clearly visible</span>
								</label>
							{/if}
						</div>

						{#if error}
							<p class="error">{error}</p>
						{/if}

						<div class="btn-group">
							<button class="btn-secondary" onclick={prevStep}>Back</button>
							<button class="btn-primary" onclick={nextStep} disabled={!documentUrl || uploadingDocument}>
								Continue
							</button>
						</div>
					</div>

				{:else if step === 3}
					<!-- Step 3: Video Recording -->
					<div class="step-content">
						<h2>Record Verification Video</h2>
						<p class="step-description">Record a short video saying your name and the code below:</p>

						<div class="verification-code-box">
							<span class="code-label">Your verification code:</span>
							<span class="code">{verificationCode}</span>
						</div>

						<div class="video-instructions">
							<h4>Please say in your video:</h4>
							<p class="script">"My name is [Your Name] and my verification code is <strong>{verificationCode}</strong>"</p>
						</div>

						<div class="upload-area">
							{#if videoPreview}
								<div class="preview-container video">
									<video src={videoPreview} controls class="video-preview"></video>
									{#if uploadingVideo}
										<div class="upload-overlay">
											<div class="spinner"></div>
											<span>Uploading...</span>
										</div>
									{:else}
										<button class="remove-btn" onclick={() => { videoFile = null; videoPreview = null; videoUrl = null; }}>
											Re-record
										</button>
									{/if}
								</div>
							{:else}
								<label class="upload-box">
									<input type="file" accept="video/*" onchange={handleVideoChange} />
									<span class="upload-icon">🎥</span>
									<span class="upload-text">Click to upload video</span>
									<span class="upload-hint">Max 50MB • MP4, MOV, or WebM</span>
								</label>
							{/if}
						</div>

						{#if error}
							<p class="error">{error}</p>
						{/if}

						<div class="btn-group">
							<button class="btn-secondary" onclick={prevStep}>Back</button>
							<button class="btn-primary" onclick={nextStep} disabled={!videoUrl || uploadingVideo}>
								Continue
							</button>
						</div>
					</div>

				{:else if step === 4}
					<!-- Step 4: Review & Submit -->
					<div class="step-content">
						<h2>Review & Submit</h2>
						<p class="step-description">Please review your submission before submitting:</p>

						<div class="review-section">
							<div class="review-item">
								<span class="review-label">Document Type</span>
								<span class="review-value">{documentTypes.find(d => d.value === documentType)?.label}</span>
							</div>
							<div class="review-item">
								<span class="review-label">Document</span>
								<span class="review-value">✓ Uploaded</span>
							</div>
							<div class="review-item">
								<span class="review-label">Video</span>
								<span class="review-value">✓ Uploaded</span>
							</div>
							<div class="review-item">
								<span class="review-label">Verification Code</span>
								<span class="review-value code">{verificationCode}</span>
							</div>
						</div>

						<div class="privacy-notice small">
							<span class="lock-icon">🔐</span>
							<p>Your documents and video will be <strong>permanently deleted</strong> after verification.</p>
						</div>

						{#if error}
							<p class="error">{error}</p>
						{/if}

						<div class="btn-group">
							<button class="btn-secondary" onclick={prevStep}>Back</button>
							<button class="btn-primary" onclick={submitVerification} disabled={submitting}>
								{submitting ? 'Submitting...' : 'Submit for Verification'}
							</button>
						</div>
					</div>
				{/if}
			</div>
		{/if}
	</main>

	<div class="desktop-footer">
		<Footer />
	</div>
</div>

<style>
	.verify-page {
		min-height: 100vh;
		display: flex;
		flex-direction: column;
		background: #0a0a0a;
	}

	.container {
		flex: 1;
		max-width: 600px;
		margin: 0 auto;
		padding: 2rem 1.5rem;
		width: 100%;
	}

	.loading {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding: 4rem;
		color: rgba(255, 255, 255, 0.6);
	}

	.spinner {
		width: 40px;
		height: 40px;
		border: 3px solid rgba(255, 255, 255, 0.1);
		border-top-color: #fff;
		border-radius: 50%;
		animation: spin 1s linear infinite;
		margin-bottom: 1rem;
	}

	@keyframes spin {
		to { transform: rotate(360deg); }
	}

	/* Progress Bar */
	.progress-bar {
		display: flex;
		align-items: center;
		justify-content: center;
		margin-bottom: 2rem;
	}

	.step {
		width: 36px;
		height: 36px;
		border-radius: 50%;
		background: rgba(255, 255, 255, 0.1);
		display: flex;
		align-items: center;
		justify-content: center;
		transition: all 0.3s ease;
	}

	.step.active {
		background: #fff;
	}

	.step.active .step-number {
		color: #0a0a0a;
	}

	.step.current {
		box-shadow: 0 0 0 4px rgba(255, 255, 255, 0.2);
	}

	.step-number {
		color: rgba(255, 255, 255, 0.5);
		font-weight: 600;
		font-size: 0.9rem;
	}

	.step-line {
		width: 40px;
		height: 2px;
		background: rgba(255, 255, 255, 0.1);
		margin: 0 0.5rem;
		transition: background 0.3s ease;
	}

	.step-line.active {
		background: #fff;
	}

	/* Step Content */
	.step-content {
		background: rgba(255, 255, 255, 0.03);
		border: 1px solid rgba(255, 255, 255, 0.1);
		padding: 2rem;
	}

	.header-section {
		text-align: center;
		margin-bottom: 2rem;
	}

	.shield-icon {
		font-size: 3rem;
		margin-bottom: 1rem;
	}

	h1 {
		font-size: 1.75rem;
		margin-bottom: 0.5rem;
	}

	h2 {
		font-size: 1.25rem;
		margin-bottom: 0.5rem;
	}

	.subtitle, .step-description {
		color: rgba(255, 255, 255, 0.6);
		font-size: 0.9rem;
		margin-bottom: 1.5rem;
	}

	/* Benefits Grid */
	.benefits-grid {
		display: grid;
		grid-template-columns: repeat(2, 1fr);
		gap: 1rem;
		margin-bottom: 2rem;
	}

	.benefit-card {
		background: rgba(255, 255, 255, 0.05);
		border: 1px solid rgba(255, 255, 255, 0.15);
		padding: 1rem;
		text-align: center;
	}

	.benefit-icon {
		font-size: 1.5rem;
		display: block;
		margin-bottom: 0.5rem;
	}

	.benefit-card h3 {
		font-size: 0.9rem;
		margin-bottom: 0.25rem;
		color: #fff;
	}

	.benefit-card p {
		font-size: 0.75rem;
		color: rgba(255, 255, 255, 0.6);
		margin: 0;
	}

	/* Privacy Notice */
	.privacy-notice {
		display: flex;
		gap: 1rem;
		background: rgba(255, 255, 255, 0.05);
		padding: 1rem;
		margin-bottom: 1.5rem;
		align-items: flex-start;
	}

	.privacy-notice.small {
		padding: 0.75rem;
		font-size: 0.85rem;
	}

	.privacy-notice .lock-icon {
		font-size: 1.25rem;
		flex-shrink: 0;
	}

	.privacy-notice strong {
		color: #fff;
	}

	.privacy-notice p {
		color: rgba(255, 255, 255, 0.7);
		font-size: 0.85rem;
		margin: 0.25rem 0 0;
	}

	/* Document Types */
	.document-types {
		display: flex;
		gap: 0.75rem;
		margin-bottom: 1.5rem;
	}

	.document-option {
		flex: 1;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.5rem;
		padding: 1rem;
		background: rgba(255, 255, 255, 0.03);
		border: 1px solid rgba(255, 255, 255, 0.1);
		cursor: pointer;
		transition: all 0.2s ease;
	}

	.document-option:hover {
		background: rgba(255, 255, 255, 0.05);
	}

	.document-option.selected {
		background: rgba(255, 255, 255, 0.1);
		border-color: #fff;
	}

	.document-option input {
		display: none;
	}

	.doc-icon {
		font-size: 1.5rem;
	}

	.doc-label {
		font-size: 0.75rem;
		text-align: center;
		color: rgba(255, 255, 255, 0.8);
	}

	/* Upload Area */
	.upload-area {
		margin-bottom: 1.5rem;
	}

	.upload-box {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding: 3rem 2rem;
		border: 2px dashed rgba(255, 255, 255, 0.2);
		cursor: pointer;
		transition: all 0.2s ease;
	}

	.upload-box:hover {
		border-color: rgba(255, 255, 255, 0.5);
		background: rgba(255, 255, 255, 0.05);
	}

	.upload-box input {
		display: none;
	}

	.upload-icon {
		font-size: 2.5rem;
		margin-bottom: 1rem;
	}

	.upload-text {
		font-size: 0.9rem;
		color: #fff;
		margin-bottom: 0.25rem;
	}

	.upload-hint {
		font-size: 0.75rem;
		color: rgba(255, 255, 255, 0.5);
	}

	.preview-container {
		position: relative;
	}

	.document-preview {
		width: 100%;
		max-height: 300px;
		object-fit: contain;
		background: #111;
	}

	.video-preview {
		width: 100%;
		max-height: 300px;
		background: #111;
	}

	.upload-overlay {
		position: absolute;
		inset: 0;
		background: rgba(0, 0, 0, 0.7);
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		color: #fff;
		gap: 0.5rem;
	}

	.upload-success {
		position: absolute;
		bottom: 0.5rem;
		left: 0.5rem;
		background: rgba(255, 255, 255, 0.9);
		color: #0a0a0a;
		padding: 0.25rem 0.75rem;
		font-size: 0.8rem;
		font-weight: 600;
	}

	.upload-failed {
		position: absolute;
		bottom: 0.5rem;
		left: 0.5rem;
		background: rgba(239, 68, 68, 0.9);
		color: #fff;
		padding: 0.25rem 0.75rem;
		font-size: 0.8rem;
		font-weight: 600;
	}

	.remove-btn {
		position: absolute;
		top: 0.5rem;
		right: 0.5rem;
		background: rgba(0, 0, 0, 0.7);
		border: none;
		color: #fff;
		padding: 0.5rem 1rem;
		font-size: 0.8rem;
		cursor: pointer;
	}

	/* Verification Code Box */
	.verification-code-box {
		background: rgba(255, 255, 255, 0.1);
		border: 2px solid #fff;
		padding: 1.5rem;
		text-align: center;
		margin-bottom: 1.5rem;
	}

	.code-label {
		display: block;
		font-size: 0.85rem;
		color: rgba(255, 255, 255, 0.7);
		margin-bottom: 0.5rem;
	}

	.code {
		font-size: 2rem;
		font-weight: 700;
		letter-spacing: 0.5rem;
		color: #fff;
		font-family: 'Courier New', monospace;
	}

	.video-instructions {
		background: rgba(255, 255, 255, 0.05);
		padding: 1rem;
		margin-bottom: 1.5rem;
	}

	.video-instructions h4 {
		font-size: 0.85rem;
		margin-bottom: 0.5rem;
		color: rgba(255, 255, 255, 0.8);
	}

	.script {
		font-size: 0.9rem;
		color: #fff;
		font-style: italic;
		margin: 0;
	}

	/* Review Section */
	.review-section {
		background: rgba(255, 255, 255, 0.03);
		padding: 1rem;
		margin-bottom: 1.5rem;
	}

	.review-item {
		display: flex;
		justify-content: space-between;
		padding: 0.75rem 0;
		border-bottom: 1px solid rgba(255, 255, 255, 0.1);
	}

	.review-item:last-child {
		border-bottom: none;
	}

	.review-label {
		color: rgba(255, 255, 255, 0.6);
		font-size: 0.9rem;
	}

	.review-value {
		color: #fff;
		font-size: 0.9rem;
	}

	.review-value.code {
		font-family: 'Courier New', monospace;
		letter-spacing: 0.1rem;
	}

	/* Buttons */
	.btn-group {
		display: flex;
		gap: 1rem;
	}

	.btn-primary, .btn-secondary {
		flex: 1;
		padding: 0.875rem 1.5rem;
		font-size: 0.9rem;
		font-weight: 600;
		font-family: 'Montserrat', sans-serif;
		cursor: pointer;
		transition: all 0.2s ease;
		border: none;
	}

	.btn-primary {
		background: #fff;
		color: #0a0a0a;
	}

	.btn-primary:hover:not(:disabled) {
		transform: translateY(-1px);
		box-shadow: 0 4px 12px rgba(255, 255, 255, 0.2);
	}

	.btn-primary:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.btn-secondary {
		background: transparent;
		border: 1px solid rgba(255, 255, 255, 0.2);
		color: rgba(255, 255, 255, 0.8);
	}

	.btn-secondary:hover {
		background: rgba(255, 255, 255, 0.05);
	}

	.error {
		color: #ef4444;
		font-size: 0.85rem;
		margin-bottom: 1rem;
		text-align: center;
	}

	/* Status Card */
	.status-card {
		text-align: center;
		padding: 3rem 2rem;
		background: rgba(255, 255, 255, 0.03);
		border: 1px solid rgba(255, 255, 255, 0.1);
	}

	.status-icon {
		font-size: 4rem;
		margin-bottom: 1.5rem;
	}

	.status-icon.pending {
		animation: pulse 2s ease-in-out infinite;
	}

	@keyframes pulse {
		0%, 100% { transform: scale(1); }
		50% { transform: scale(1.1); }
	}

	.status-card h1 {
		margin-bottom: 1rem;
	}

	.status-card p {
		color: rgba(255, 255, 255, 0.7);
		margin-bottom: 0.5rem;
	}

	.status-card .hint {
		font-size: 0.85rem;
		color: rgba(255, 255, 255, 0.5);
		margin-bottom: 2rem;
	}

	.status-card .btn-primary {
		display: inline-block;
		text-decoration: none;
		padding: 0.875rem 2rem;
	}

	.mobile-close {
		display: none;
	}

	@media (max-width: 768px) {
		.desktop-header,
		.desktop-footer {
			display: none;
		}

		.mobile-close {
			display: flex;
			position: fixed;
			top: 1rem;
			right: 1rem;
			width: 44px;
			height: 44px;
			align-items: center;
			justify-content: center;
			color: rgba(255, 255, 255, 0.7);
			text-decoration: none;
			z-index: 10;
		}

		.mobile-close:hover {
			color: #fff;
		}

		.verify-page {
			position: relative;
			min-height: 100vh;
			padding-bottom: 5rem;
		}

		.container {
			padding: 3.5rem 1rem 1rem;
		}

		.step-content {
			padding: 1.5rem;
			background: transparent;
			border: none;
		}

		.benefits-grid {
			grid-template-columns: 1fr;
		}

		.document-types {
			flex-direction: row;
			flex-wrap: wrap;
		}

		.document-option {
			flex: 1;
			min-width: 0;
			padding: 0.75rem 0.5rem;
		}

		.doc-label {
			font-size: 0.65rem;
		}

		.doc-icon {
			font-size: 1.25rem;
		}

		.code {
			font-size: 1.5rem;
			letter-spacing: 0.3rem;
		}

		.btn-group {
			position: fixed;
			bottom: 0;
			left: 0;
			right: 0;
			padding: 1rem;
			background: #0a0a0a;
			border-top: 1px solid rgba(255, 255, 255, 0.1);
			z-index: 10;
		}

		.btn-group.single {
			display: flex;
		}

		.btn-group.single .btn-primary {
			width: 100%;
		}

		.status-card {
			background: transparent;
			border: none;
			padding: 2rem 0;
		}

		.status-card .btn-primary {
			position: fixed;
			bottom: 1rem;
			left: 1rem;
			right: 1rem;
			width: auto;
		}
	}
</style>

