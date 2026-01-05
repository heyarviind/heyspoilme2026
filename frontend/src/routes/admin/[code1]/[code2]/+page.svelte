<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';

	const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8080';

	let code1 = $state('');
	let code2 = $state('');
	let loading = $state(true);
	let error = $state('');
	let activeTab = $state<'users' | 'messages' | 'verifications' | 'images'>('users');

	// Stats
	let stats = $state<Record<string, number>>({});

	// Users
	let users = $state<any[]>([]);
	let usersTotal = $state(0);
	let usersPage = $state(1);
	let userSearch = $state('');
	let selectedUser = $state<any>(null);

	// Messages
	let messages = $state<any[]>([]);
	let messagesTotal = $state(0);
	let messagesPage = $state(1);

	// Verifications
	let verifications = $state<any[]>([]);
	let verificationFilter = $state('pending');

	// Images
	let allImages = $state<any[]>([]);
	let imagesTotal = $state(0);
	let imagesPage = $state(1);

	// Modals
	let showUserModal = $state(false);
	let showDeleteConfirm = $state(false);
	let deleteTarget = $state<{ type: string; id: string; name?: string } | null>(null);
	let showRejectModal = $state(false);
	let rejectTarget = $state<any>(null);
	let rejectReason = $state('');

	// Image viewer
	let showImageViewer = $state(false);
	let viewerImageUrl = $state('');

	$effect(() => {
		const unsubscribe = page.subscribe(p => {
			code1 = p.params.code1 || '';
			code2 = p.params.code2 || '';
		});
		return unsubscribe;
	});

	onMount(async () => {
		const p = $page;
		code1 = p.params.code1 || '';
		code2 = p.params.code2 || '';
		await loadStats();
		await loadData();
	});

	async function adminFetch(endpoint: string, options: RequestInit = {}) {
		const response = await fetch(`${API_BASE}/admin/${code1}/${code2}${endpoint}`, {
			...options,
			headers: {
				'Content-Type': 'application/json',
				...options.headers,
			},
		});

		if (!response.ok) {
			if (response.status === 401) {
				goto('/auth/login');
				throw new Error('Unauthorized');
			}
			if (response.status === 404) {
				throw new Error('Invalid admin access');
			}
			const err = await response.json().catch(() => ({ error: 'Request failed' }));
			throw new Error(err.error || 'Request failed');
		}

		return response.json();
	}

	async function loadStats() {
		try {
			stats = await adminFetch('/stats');
		} catch (e: any) {
			error = e.message;
		}
	}

	async function loadData() {
		loading = true;
		error = '';

		try {
			if (activeTab === 'users') {
				const params = new URLSearchParams();
				params.append('page', String(usersPage));
				params.append('limit', '20');
				if (userSearch) params.append('search', userSearch);
				
				const result = await adminFetch(`/users?${params.toString()}`);
				users = result.users || [];
				usersTotal = result.total || 0;
			} else if (activeTab === 'messages') {
				const params = new URLSearchParams();
				params.append('page', String(messagesPage));
				params.append('limit', '50');
				
				const result = await adminFetch(`/messages?${params.toString()}`);
				messages = result.messages || [];
				messagesTotal = result.total || 0;
			} else if (activeTab === 'verifications') {
				const result = await adminFetch(`/verifications?status=${verificationFilter}`);
				verifications = result.requests || [];
			} else if (activeTab === 'images') {
				const params = new URLSearchParams();
				params.append('page', String(imagesPage));
				params.append('limit', '50');
				
				const result = await adminFetch(`/images?${params.toString()}`);
				allImages = result.images || [];
				imagesTotal = result.total || 0;
			}
		} catch (e: any) {
			error = e.message;
		} finally {
			loading = false;
		}
	}

	async function viewUser(userId: string) {
		try {
			selectedUser = await adminFetch(`/users/${userId}`);
			showUserModal = true;
		} catch (e: any) {
			alert('Error loading user: ' + e.message);
		}
	}

	async function deleteProfileImage(imageId: string) {
		if (!confirm('Delete this image?')) return;
		
		try {
			await adminFetch(`/images/${imageId}`, { method: 'DELETE' });
			// Refresh user
			if (selectedUser) {
				selectedUser.images = selectedUser.images.filter((img: any) => img.id !== imageId);
			}
			await loadStats();
		} catch (e: any) {
			alert('Error: ' + e.message);
		}
	}

	async function deleteImage(imageId: string, source: string) {
		if (!confirm('Delete this image?')) return;
		
		try {
			await adminFetch(`/images/${imageId}?source=${source}`, { method: 'DELETE' });
			allImages = allImages.filter((img: any) => img.id !== imageId);
		} catch (e: any) {
			alert('Error: ' + e.message);
		}
	}

	async function confirmDelete() {
		if (!deleteTarget) return;
		
		try {
			if (deleteTarget.type === 'user') {
				await adminFetch(`/users/${deleteTarget.id}`, { method: 'DELETE' });
				showUserModal = false;
				selectedUser = null;
				await loadData();
				await loadStats();
			}
		} catch (e: any) {
			alert('Error: ' + e.message);
		} finally {
			showDeleteConfirm = false;
			deleteTarget = null;
		}
	}

	async function approveVerification(requestId: string) {
		try {
			await adminFetch(`/verifications/${requestId}/approve`, { method: 'POST' });
			await loadData();
			await loadStats();
		} catch (e: any) {
			alert('Error: ' + e.message);
		}
	}

	async function rejectVerification() {
		if (!rejectTarget || !rejectReason.trim()) return;
		
		try {
			await adminFetch(`/verifications/${rejectTarget.id}/reject`, {
				method: 'POST',
				body: JSON.stringify({ reason: rejectReason }),
			});
			await loadData();
			await loadStats();
		} catch (e: any) {
			alert('Error: ' + e.message);
		} finally {
			showRejectModal = false;
			rejectTarget = null;
			rejectReason = '';
		}
	}

	async function updateWealthStatus(userId: string, status: string) {
		try {
			await adminFetch(`/users/${userId}/wealth-status`, {
				method: 'PUT',
				body: JSON.stringify({ status }),
			});
			if (selectedUser && selectedUser.id === userId) {
				selectedUser.wealth_status = status;
			}
			await loadData();
			await loadStats();
		} catch (e: any) {
			alert('Error: ' + e.message);
		}
	}

	function formatDate(dateStr: string) {
		if (!dateStr) return '-';
		return new Date(dateStr).toLocaleString();
	}

	function formatRelativeTime(dateStr: string) {
		if (!dateStr) return 'never';
		const date = new Date(dateStr);
		const now = new Date();
		const diff = now.getTime() - date.getTime();
		const mins = Math.floor(diff / 60000);
		if (mins < 1) return 'just now';
		if (mins < 60) return `${mins}m ago`;
		const hours = Math.floor(mins / 60);
		if (hours < 24) return `${hours}h ago`;
		const days = Math.floor(hours / 24);
		return `${days}d ago`;
	}

	function getStatusBadgeClass(status: string) {
		if (status === 'low') return 'badge-trusted';
		if (status === 'medium') return 'badge-premium';
		if (status === 'high') return 'badge-elite';
		return 'badge-standard';
	}

	function getStatusLabel(status: string) {
		if (status === 'low') return 'Trusted';
		if (status === 'medium') return 'Premium';
		if (status === 'high') return 'Elite';
		return 'Standard';
	}
</script>

<svelte:head>
	<title>Admin Panel - HeySpoilMe</title>
	<link href="https://fonts.googleapis.com/css2?family=Montserrat:wght@400;500;600;700&display=swap" rel="stylesheet">
</svelte:head>

<div class="admin-container">
	<header class="site-header">
		<a href="/browse" class="logo-link">
			<img src="/img/logo.svg" alt="HeySpoilMe" class="logo" />
		</a>
		
		<nav class="header-nav">
			<button 
				class="nav-link" 
				class:active={activeTab === 'users'}
				onclick={() => { activeTab = 'users'; usersPage = 1; loadData(); }}
			>
				Users
			</button>
			<button 
				class="nav-link" 
				class:active={activeTab === 'messages'}
				onclick={() => { activeTab = 'messages'; messagesPage = 1; loadData(); }}
			>
				Messages
			</button>
			<button 
				class="nav-link" 
				class:active={activeTab === 'verifications'}
				onclick={() => { activeTab = 'verifications'; loadData(); }}
			>
				Verifications
				{#if stats.pending_verifications}
					<span class="nav-badge">{stats.pending_verifications}</span>
				{/if}
			</button>
			<button 
				class="nav-link" 
				class:active={activeTab === 'images'}
				onclick={() => { activeTab = 'images'; imagesPage = 1; loadData(); }}
			>
				Images
			</button>
		</nav>
	</header>

	<div class="admin-subheader">
		<div class="toolbar">
			<input 
				type="text" 
				placeholder="Search by email or name..." 
				bind:value={userSearch}
				onkeydown={(e) => e.key === 'Enter' && loadData()}
			/>
			<button class="btn-search" onclick={() => loadData()}>Search</button>
		</div>
		<div class="header-stats">
			<div class="stat">
				<span class="stat-value">{stats.total_users ?? '-'}</span>
				<span class="stat-label">Users</span>
			</div>
			<div class="stat">
				<span class="stat-value online">{stats.online_users ?? '-'}</span>
				<span class="stat-label">Online</span>
			</div>
			<div class="stat">
				<span class="stat-value">{stats.pending_verifications ?? '-'}</span>
				<span class="stat-label">Pending</span>
			</div>
			<div class="stat">
				<span class="stat-value verified">{stats.verified_users ?? '-'}</span>
				<span class="stat-label">Verified</span>
			</div>
			<div class="stat">
				<span class="stat-value trusted">{stats.trusted_members ?? '-'}</span>
				<span class="stat-label">Trusted</span>
			</div>
		</div>
	</div>

	{#if error}
		<div class="error-banner">{error}</div>
	{/if}

	<main class="content">
		{#if loading}
			<div class="loading">Loading...</div>
		{:else if activeTab === 'users'}

			<div class="table-container">
				<table>
					<thead>
						<tr>
							<th>User</th>
							<th>Gender</th>
							<th>Email</th>
							<th>Email</th>
							<th>ID</th>
							<th>Status</th>
							<th>Online</th>
							<th>Imgs</th>
							<th>Created</th>
							<th>Actions</th>
						</tr>
					</thead>
					<tbody>
						{#each users as user}
							<tr>
								<td>
									<div class="user-cell">
										<span class="user-name">{user.display_name || 'No profile'}</span>
										{#if user.age}
											<span class="user-meta">{user.age} yrs</span>
										{/if}
									</div>
								</td>
								<td>
									{#if user.gender === 'male'}
										<span class="gender-badge male">Male</span>
									{:else if user.gender === 'female'}
										<span class="gender-badge female">Female</span>
									{:else}
										<span class="gender-badge">-</span>
									{/if}
								</td>
								<td class="email-cell">{user.email}</td>
								<td>
									{#if user.email_verified}
										<span class="status-dot verified" title="Email Verified"></span>
									{:else}
										<span class="status-dot" title="Email Not Verified"></span>
									{/if}
								</td>
								<td>
									{#if user.is_verified}
										<span class="status-dot verified" title="ID Verified"></span>
									{:else}
										<span class="status-dot" title="ID Not Verified"></span>
									{/if}
								</td>
								<td>
									<span class="badge {getStatusBadgeClass(user.wealth_status)}">
										{getStatusLabel(user.wealth_status)}
									</span>
								</td>
								<td>
									{#if user.is_online}
										<span class="status-dot online" title="Online"></span>
									{:else}
										<span class="offline-time">{formatRelativeTime(user.last_seen)}</span>
									{/if}
								</td>
								<td>{user.image_count}</td>
								<td class="date-cell">{formatDate(user.created_at)}</td>
								<td>
									<button class="btn-small" onclick={() => viewUser(user.id)}>View</button>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>

			<div class="pagination">
				<button 
					disabled={usersPage <= 1} 
					onclick={() => { usersPage--; loadData(); }}
				>← Prev</button>
				<span>Page {usersPage} of {Math.ceil(usersTotal / 20)}</span>
				<button 
					disabled={usersPage >= Math.ceil(usersTotal / 20)} 
					onclick={() => { usersPage++; loadData(); }}
				>Next →</button>
			</div>

		{:else if activeTab === 'messages'}
			<div class="table-container">
				<table>
					<thead>
						<tr>
							<th>Sender</th>
							<th>Email</th>
							<th>Content</th>
							<th>Sent</th>
						</tr>
					</thead>
					<tbody>
						{#each messages as msg}
							<tr>
								<td>{msg.sender_name || 'Unknown'}</td>
								<td class="email-cell">{msg.sender_email}</td>
								<td class="content-cell">{msg.content}</td>
								<td class="date-cell">{formatDate(msg.created_at)}</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>

			<div class="pagination">
				<button 
					disabled={messagesPage <= 1} 
					onclick={() => { messagesPage--; loadData(); }}
				>← Prev</button>
				<span>Page {messagesPage} of {Math.ceil(messagesTotal / 50)}</span>
				<button 
					disabled={messagesPage >= Math.ceil(messagesTotal / 50)} 
					onclick={() => { messagesPage++; loadData(); }}
				>Next →</button>
			</div>

		{:else if activeTab === 'verifications'}
			<div class="toolbar">
				<select bind:value={verificationFilter} onchange={() => loadData()}>
					<option value="pending">Pending</option>
					<option value="approved">Approved</option>
					<option value="rejected">Rejected</option>
					<option value="">All</option>
				</select>
			</div>

			<div class="verifications-grid">
				{#each verifications as req}
					<div class="verification-card">
						<div class="verification-header">
							<div>
								<span class="user-name">{req.display_name || 'Unknown'}</span>
								<span class="user-email">{req.user_email}</span>
							</div>
							<span class="verification-status {req.status}">{req.status}</span>
						</div>
						
						<div class="verification-info">
							<div><strong>Document:</strong> {req.document_type}</div>
							<div><strong>Code:</strong> {req.verification_code}</div>
							<div><strong>Submitted:</strong> {formatDate(req.created_at)}</div>
							{#if req.rejection_reason}
								<div class="rejection-reason"><strong>Reason:</strong> {req.rejection_reason}</div>
							{/if}
						</div>

						<div class="verification-media">
							<button 
								class="media-btn" 
								onclick={() => { viewerImageUrl = req.document_url; showImageViewer = true; }}
							>
								Document
							</button>
							<button 
								class="media-btn" 
								onclick={() => { viewerImageUrl = req.video_url; showImageViewer = true; }}
							>
								Video
							</button>
						</div>

						{#if req.status === 'pending'}
							<div class="verification-actions">
								<button class="btn-approve" onclick={() => approveVerification(req.id)}>
									Approve
								</button>
								<button class="btn-reject" onclick={() => { rejectTarget = req; showRejectModal = true; }}>
									Reject
								</button>
							</div>
						{/if}
					</div>
				{/each}

				{#if verifications.length === 0}
					<div class="empty-state">No verification requests found.</div>
				{/if}
			</div>

		{:else if activeTab === 'images'}
			<div class="images-gallery">
				{#each allImages as img}
					<div class="gallery-card">
						<div class="gallery-image" onclick={() => { viewerImageUrl = img.url; showImageViewer = true; }}>
							{#if img.url.includes('.mp4') || img.url.includes('.webm') || img.url.includes('.mov')}
								<video src={img.url}></video>
							{:else}
								<img src={img.url} alt="Uploaded" />
							{/if}
						</div>
						<div class="gallery-info">
							<div class="gallery-type">{img.source}</div>
							<div class="gallery-user">{img.user_email}</div>
							<div class="gallery-date">{formatDate(img.created_at)}</div>
						</div>
						<button class="gallery-delete" onclick={() => deleteImage(img.id, img.source)}>Delete</button>
					</div>
				{/each}

				{#if allImages.length === 0}
					<div class="empty-state">No images found.</div>
				{/if}
			</div>

			<div class="pagination">
				<button 
					disabled={imagesPage <= 1} 
					onclick={() => { imagesPage--; loadData(); }}
				>Prev</button>
				<span>Page {imagesPage} of {Math.ceil(imagesTotal / 50) || 1}</span>
				<button 
					disabled={imagesPage >= Math.ceil(imagesTotal / 50)} 
					onclick={() => { imagesPage++; loadData(); }}
				>Next</button>
			</div>
		{/if}
	</main>
</div>

<!-- User Detail Modal -->
{#if showUserModal && selectedUser}
	<div class="modal-overlay" onclick={() => showUserModal = false}>
		<div class="modal" onclick={(e) => e.stopPropagation()}>
			<div class="modal-header">
				<h2>{selectedUser.display_name || selectedUser.email}</h2>
				<button class="modal-close" onclick={() => showUserModal = false}>X</button>
			</div>
			
			<div class="modal-body">
				<div class="user-details">
					<div class="detail-row">
						<span class="label">Email:</span>
						<span class="value">{selectedUser.email}</span>
					</div>
					<div class="detail-row">
						<span class="label">Email Verified:</span>
						<span class="value">{selectedUser.email_verified ? 'Yes' : 'No'}</span>
					</div>
					<div class="detail-row">
						<span class="label">ID Verified:</span>
						<span class="value">{selectedUser.is_verified ? 'Yes' : 'No'}</span>
					</div>
					<div class="detail-row">
						<span class="label">Gender:</span>
						<span class="value">{selectedUser.gender || '-'}</span>
					</div>
					<div class="detail-row">
						<span class="label">Age:</span>
						<span class="value">{selectedUser.age || '-'}</span>
					</div>
					<div class="detail-row">
						<span class="label">Location:</span>
						<span class="value">{selectedUser.city ? `${selectedUser.city}, ${selectedUser.state}` : '-'}</span>
					</div>
					<div class="detail-row">
						<span class="label">Status:</span>
						<span class="value">
							<select 
								value={selectedUser.wealth_status}
								onchange={(e) => updateWealthStatus(selectedUser.id, e.currentTarget.value)}
							>
								<option value="none">Standard</option>
								<option value="low">Trusted</option>
								<option value="medium">Premium</option>
								<option value="high">Elite</option>
							</select>
						</span>
					</div>
					<div class="detail-row">
						<span class="label">Online:</span>
						<span class="value">
							{#if selectedUser.is_online}
								<span class="online-dot"></span> Online
							{:else}
								Last seen {formatRelativeTime(selectedUser.last_seen)}
							{/if}
						</span>
					</div>
					<div class="detail-row">
						<span class="label">Joined:</span>
						<span class="value">{formatDate(selectedUser.created_at)}</span>
					</div>
				</div>

				{#if selectedUser.images && selectedUser.images.length > 0}
					<h3>Profile Images</h3>
					<div class="images-grid">
						{#each selectedUser.images as img}
							<div class="image-card">
								<img 
									src={img.url} 
									alt="Profile" 
									onclick={() => { viewerImageUrl = img.url; showImageViewer = true; }}
								/>
								<button class="delete-image" onclick={() => deleteProfileImage(img.id)}>
									Delete
								</button>
							</div>
						{/each}
					</div>
				{:else}
					<p class="no-images">No profile images</p>
				{/if}
			</div>

			<div class="modal-footer">
				<button 
					class="btn-danger" 
					onclick={() => { 
						deleteTarget = { type: 'user', id: selectedUser.id, name: selectedUser.display_name || selectedUser.email };
						showDeleteConfirm = true;
					}}
				>
					Delete User
				</button>
			</div>
		</div>
	</div>
{/if}

<!-- Delete Confirmation Modal -->
{#if showDeleteConfirm && deleteTarget}
	<div class="modal-overlay" onclick={() => showDeleteConfirm = false}>
		<div class="modal modal-small" onclick={(e) => e.stopPropagation()}>
			<div class="modal-header">
				<h2>Confirm Delete</h2>
			</div>
			<div class="modal-body">
				<p>Are you sure you want to delete <strong>{deleteTarget.name}</strong>?</p>
				<p class="warning">This will delete all their data including messages, likes, and images.</p>
			</div>
			<div class="modal-footer">
				<button class="btn-secondary" onclick={() => showDeleteConfirm = false}>Cancel</button>
				<button class="btn-danger" onclick={confirmDelete}>Delete</button>
			</div>
		</div>
	</div>
{/if}

<!-- Reject Modal -->
{#if showRejectModal && rejectTarget}
	<div class="modal-overlay" onclick={() => showRejectModal = false}>
		<div class="modal modal-small" onclick={(e) => e.stopPropagation()}>
			<div class="modal-header">
				<h2>Reject Verification</h2>
			</div>
			<div class="modal-body">
				<p>Rejecting verification for: <strong>{rejectTarget.display_name || rejectTarget.user_email}</strong></p>
				<textarea 
					placeholder="Reason for rejection..." 
					bind:value={rejectReason}
					rows="3"
				></textarea>
			</div>
			<div class="modal-footer">
				<button class="btn-secondary" onclick={() => showRejectModal = false}>Cancel</button>
				<button class="btn-danger" onclick={rejectVerification} disabled={!rejectReason.trim()}>Reject</button>
			</div>
		</div>
	</div>
{/if}

<!-- Image Viewer Modal -->
{#if showImageViewer}
	<div class="modal-overlay image-viewer" onclick={() => showImageViewer = false}>
		<button class="close-viewer" onclick={() => showImageViewer = false}>X</button>
		{#if viewerImageUrl.includes('.mp4') || viewerImageUrl.includes('.webm') || viewerImageUrl.includes('.mov')}
			<video src={viewerImageUrl} controls autoplay></video>
		{:else}
			<img src={viewerImageUrl} alt="Full size" />
		{/if}
	</div>
{/if}

<style>
	.admin-container {
		min-height: 100vh;
		background: #0a0a0a;
		color: #fff;
		font-family: 'Montserrat', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
	}

	.site-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 1rem 2rem;
		border-bottom: 1px solid rgba(255, 255, 255, 0.1);
		position: sticky;
		top: 0;
		background: #0a0a0a;
		z-index: 100;
	}

	.logo-link {
		text-decoration: none;
	}

	.logo {
		height: 2.5rem;
	}

	.header-nav {
		display: flex;
		align-items: center;
		gap: 1.5rem;
	}

	.nav-link {
		background: none;
		border: none;
		color: rgba(255, 255, 255, 0.6);
		font-family: inherit;
		font-size: 0.9rem;
		font-weight: 500;
		cursor: pointer;
		padding: 0;
		transition: color 0.2s ease;
		display: flex;
		align-items: center;
		gap: 0.35rem;
	}

	.nav-link:hover,
	.nav-link.active {
		color: #fff;
	}

	.nav-badge {
		background: #ef4444;
		color: #fff;
		font-size: 0.7rem;
		font-weight: 600;
		padding: 0.15rem 0.4rem;
		min-width: 1.25rem;
		text-align: center;
	}

	.admin-subheader {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 1rem 2rem;
		background: #0a0a0a;
		gap: 2rem;
	}

	.header-stats {
		display: flex;
		gap: 2rem;
	}

	.stat {
		display: flex;
		flex-direction: column;
		align-items: center;
	}

	.stat-value {
		font-size: 1.5rem;
		font-weight: 600;
		color: #fff;
	}

	.stat-value.online {
		color: #fff;
	}

	.stat-value.verified {
		color: #60a5fa;
	}

	.stat-value.trusted {
		color: #fbbf24;
	}

	.stat-label {
		font-size: 0.75rem;
		color: #9ca3af;
		text-transform: uppercase;
	}

	.error-banner {
		padding: 1rem 2rem;
		background: rgba(220, 38, 38, 0.1);
		color: #f87171;
	}


	.content {
		padding: 1.5rem 2rem;
	}

	.loading {
		text-align: center;
		padding: 3rem;
		color: #9ca3af;
	}

	.toolbar {
		display: flex;
		gap: 0.5rem;
		flex: 1;
		max-width: 500px;
	}

	.toolbar input,
	.toolbar select {
		padding: 0.5rem 1rem;
		background: #1a1a1a;
		border: none;
		border-radius: 0;
		color: #fff;
		font-family: inherit;
		font-size: 0.875rem;
	}

	.toolbar input {
		flex: 1;
	}

	.toolbar input:focus,
	.toolbar select:focus {
		outline: none;
		box-shadow: 0 0 0 2px rgba(255, 255, 255, 0.1);
	}

	.btn-search {
		padding: 0.5rem 1rem;
		background: #fff;
		border: none;
		border-radius: 0;
		color: #0a0a0a;
		cursor: pointer;
		font-family: inherit;
	}

	.btn-search:hover {
		background: #e5e5e5;
	}

	.table-container {
		overflow-x: auto;
		background: #0a0a0a;
		border-radius: 0;
	}

	table {
		width: 100%;
		border-collapse: collapse;
	}

	th, td {
		padding: 0.75rem 1rem;
		text-align: left;
	}

	th {
		background: #1a1a1a;
		font-size: 0.75rem;
		text-transform: uppercase;
		color: #9ca3af;
		font-weight: 600;
	}

	tr:hover {
		background: #1a1a1a;
	}

	.user-cell {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	.user-name {
		font-weight: 500;
		color: #fff;
	}

	.user-meta {
		font-size: 0.75rem;
		color: #9ca3af;
	}

	.email-cell {
		font-size: 0.875rem;
		color: #9ca3af;
	}

	.date-cell {
		font-size: 0.75rem;
		color: #9ca3af;
	}

	.content-cell {
		max-width: 300px;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.badge {
		padding: 0.25rem 0.5rem;
		border-radius: 0;
		font-size: 0.75rem;
		font-weight: 500;
	}

	.badge-standard {
		background: #262626;
		color: #9ca3af;
	}

	.badge-trusted {
		background: rgba(251, 191, 36, 0.15);
		color: #fbbf24;
	}

	.badge-premium {
		background: rgba(167, 139, 250, 0.15);
		color: #a78bfa;
	}

	.badge-elite {
		background: rgba(253, 224, 71, 0.15);
		color: #fde047;
	}

	.gender-badge {
		padding: 0.25rem 0.5rem;
		border-radius: 0;
		font-size: 0.75rem;
		font-weight: 500;
		background: #262626;
		color: #9ca3af;
	}

	.gender-badge.male {
		background: rgba(96, 165, 250, 0.15);
		color: #60a5fa;
	}

	.gender-badge.female {
		background: rgba(244, 114, 182, 0.15);
		color: #f472b6;
	}

	.verified-badge {
		color: #fff;
		font-weight: 500;
	}

	.not-verified {
		color: #6b7280;
	}

	.status-dot {
		display: inline-block;
		width: 10px;
		height: 10px;
		background: #374151;
		border-radius: 50%;
	}

	.status-dot.verified {
		background: #22c55e;
	}

	.status-dot.online {
		background: #22c55e;
	}

	.online-dot {
		display: inline-block;
		width: 8px;
		height: 8px;
		background: #22c55e;
		border-radius: 50%;
	}

	.offline-time {
		font-size: 0.75rem;
		color: #9ca3af;
	}

	.btn-small {
		padding: 0.25rem 0.75rem;
		background: #fff;
		border: none;
		border-radius: 0;
		color: #0a0a0a;
		cursor: pointer;
		font-family: inherit;
		font-size: 0.75rem;
	}

	.btn-small:hover {
		background: #e5e5e5;
	}

	.pagination {
		display: flex;
		justify-content: center;
		align-items: center;
		gap: 1rem;
		padding: 1rem;
	}

	.pagination button {
		padding: 0.5rem 1rem;
		background: #fff;
		border: none;
		border-radius: 0;
		color: #0a0a0a;
		cursor: pointer;
		font-family: inherit;
	}

	.pagination button:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.pagination span {
		color: #9ca3af;
		font-size: 0.875rem;
	}

	/* Verifications Grid */
	.verifications-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
		gap: 1rem;
	}

	.verification-card {
		background: #1a1a1a;
		border: none;
		border-radius: 0;
		padding: 1rem;
	}

	.verification-header {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		margin-bottom: 1rem;
	}

	.verification-header .user-email {
		display: block;
		font-size: 0.75rem;
		color: #9ca3af;
	}

	.verification-status {
		padding: 0.25rem 0.5rem;
		border-radius: 0;
		font-size: 0.75rem;
		font-weight: 500;
		text-transform: uppercase;
	}

	.verification-status.pending {
		background: rgba(251, 191, 36, 0.15);
		color: #fbbf24;
	}

	.verification-status.approved {
		background: rgba(255, 255, 255, 0.1);
		color: #fff;
	}

	.verification-status.rejected {
		background: rgba(248, 113, 113, 0.15);
		color: #f87171;
	}

	.verification-info {
		font-size: 0.875rem;
		margin-bottom: 1rem;
	}

	.verification-info > div {
		margin-bottom: 0.25rem;
	}

	.rejection-reason {
		color: #f87171;
	}

	.verification-media {
		display: flex;
		gap: 0.5rem;
		margin-bottom: 1rem;
	}

	.media-btn {
		flex: 1;
		padding: 0.5rem;
		background: #0a0a0a;
		border: 1px solid #fff;
		border-radius: 0;
		color: #fff;
		cursor: pointer;
		font-family: inherit;
		font-size: 0.875rem;
	}

	.media-btn:hover {
		background: #262626;
	}

	.verification-actions {
		display: flex;
		gap: 0.5rem;
	}

	.btn-approve {
		flex: 1;
		padding: 0.5rem;
		background: #fff;
		border: none;
		border-radius: 0;
		color: #0a0a0a;
		cursor: pointer;
		font-family: inherit;
	}

	.btn-approve:hover {
		background: #e5e5e5;
	}

	.btn-reject {
		flex: 1;
		padding: 0.5rem;
		background: #0a0a0a;
		border: 1px solid #f87171;
		border-radius: 0;
		color: #f87171;
		cursor: pointer;
		font-family: inherit;
	}

	.btn-reject:hover {
		background: rgba(248, 113, 113, 0.1);
	}

	.empty-state {
		grid-column: 1 / -1;
		text-align: center;
		padding: 3rem;
		color: #9ca3af;
	}

	/* Modal */
	.modal-overlay {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.8);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1000;
		padding: 1rem;
	}

	.modal {
		background: #1a1a1a;
		border: none;
		border-radius: 0;
		width: 100%;
		max-width: 600px;
		max-height: 90vh;
		overflow: hidden;
		display: flex;
		flex-direction: column;
		box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
	}

	.modal-small {
		max-width: 400px;
	}

	.modal-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 1rem 1.5rem;
	}

	.modal-header h2 {
		font-size: 1.125rem;
		color: #fff;
	}

	.modal-close {
		background: transparent;
		border: none;
		color: #9ca3af;
		font-size: 1.5rem;
		cursor: pointer;
		padding: 0;
		line-height: 1;
	}

	.modal-close:hover {
		color: #fff;
	}

	.modal-body {
		padding: 1.5rem;
		overflow-y: auto;
	}

	.modal-body h3 {
		font-size: 0.875rem;
		color: #9ca3af;
		margin: 1.5rem 0 1rem;
		text-transform: uppercase;
	}

	.modal-body p {
		margin-bottom: 0.5rem;
	}

	.modal-body .warning {
		color: #f87171;
		font-size: 0.875rem;
	}

	.modal-body textarea {
		width: 100%;
		padding: 0.75rem;
		background: #262626;
		border: none;
		border-radius: 0;
		color: #fff;
		font-family: inherit;
		font-size: 0.875rem;
		resize: vertical;
		margin-top: 0.5rem;
	}

	.modal-body textarea:focus {
		outline: none;
		box-shadow: 0 0 0 2px rgba(255, 255, 255, 0.1);
	}

	.modal-footer {
		display: flex;
		justify-content: flex-end;
		gap: 0.5rem;
		padding: 1rem 1.5rem;
	}

	.btn-secondary {
		padding: 0.5rem 1rem;
		background: #0a0a0a;
		border: 1px solid #fff;
		border-radius: 0;
		color: #fff;
		cursor: pointer;
		font-family: inherit;
	}

	.btn-secondary:hover {
		background: #262626;
	}

	.btn-danger {
		padding: 0.5rem 1rem;
		background: #fff;
		border: none;
		border-radius: 0;
		color: #0a0a0a;
		cursor: pointer;
		font-family: inherit;
	}

	.btn-danger:hover {
		background: #e5e5e5;
	}

	.btn-danger:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.user-details {
		display: grid;
		gap: 0.5rem;
	}

	.detail-row {
		display: flex;
		gap: 1rem;
	}

	.detail-row .label {
		width: 120px;
		color: #9ca3af;
		flex-shrink: 0;
	}

	.detail-row .value {
		color: #fff;
	}

	.detail-row select {
		padding: 0.25rem 0.5rem;
		background: #262626;
		border: none;
		border-radius: 0;
		color: #fff;
		font-family: inherit;
		font-size: 0.875rem;
	}

	.images-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
		gap: 0.75rem;
	}

	.image-card {
		position: relative;
		aspect-ratio: 1;
		border-radius: 0;
		overflow: hidden;
		background: #262626;
	}

	.image-card img {
		width: 100%;
		height: 100%;
		object-fit: cover;
		cursor: pointer;
	}

	.image-card .delete-image {
		position: absolute;
		top: 4px;
		right: 4px;
		width: 28px;
		height: 28px;
		background: rgba(248, 113, 113, 0.9);
		border: none;
		border-radius: 0;
		cursor: pointer;
		font-size: 0.875rem;
		color: #fff;
		opacity: 0;
		transition: opacity 0.2s;
	}

	.image-card:hover .delete-image {
		opacity: 1;
	}

	.no-images {
		color: #9ca3af;
		font-style: italic;
	}

	/* Image Viewer */
	.image-viewer {
		background: rgba(0, 0, 0, 0.95);
	}

	.image-viewer img,
	.image-viewer video {
		max-width: 90vw;
		max-height: 90vh;
		object-fit: contain;
	}

	.close-viewer {
		position: absolute;
		top: 1rem;
		right: 1rem;
		width: 40px;
		height: 40px;
		background: rgba(255, 255, 255, 0.1);
		border: none;
		border-radius: 50%;
		color: #fff;
		font-size: 1.5rem;
		cursor: pointer;
	}

	.close-viewer:hover {
		background: rgba(255, 255, 255, 0.2);
	}

	/* Images Gallery */
	.images-gallery {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
		gap: 1rem;
	}

	.gallery-card {
		background: #1a1a1a;
		overflow: hidden;
	}

	.gallery-image {
		aspect-ratio: 1;
		cursor: pointer;
		overflow: hidden;
	}

	.gallery-image img,
	.gallery-image video {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	.gallery-info {
		padding: 0.75rem;
	}

	.gallery-type {
		font-size: 0.75rem;
		color: #9ca3af;
		text-transform: uppercase;
		margin-bottom: 0.25rem;
	}

	.gallery-user {
		font-size: 0.75rem;
		color: #fff;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.gallery-date {
		font-size: 0.625rem;
		color: #6b7280;
		margin-top: 0.25rem;
	}

	.gallery-delete {
		width: 100%;
		padding: 0.5rem;
		background: transparent;
		border: none;
		border-top: 1px solid #262626;
		color: #f87171;
		cursor: pointer;
		font-family: inherit;
		font-size: 0.75rem;
	}

	.gallery-delete:hover {
		background: rgba(248, 113, 113, 0.1);
	}

	/* Responsive */
	@media (max-width: 768px) {
		.admin-header {
			flex-direction: column;
			gap: 1rem;
		}

		.header-stats {
			flex-wrap: wrap;
			justify-content: center;
		}

		.header-nav {
			gap: 1rem;
		}

		.content {
			padding: 1rem;
		}

		.verifications-grid {
			grid-template-columns: 1fr;
		}

		.images-gallery {
			grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
		}
	}
</style>

