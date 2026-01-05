-- Create verification_requests table for identity verification
CREATE TABLE verification_requests (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    document_type VARCHAR(50) NOT NULL, -- 'aadhar', 'passport', 'driving_license'
    document_url TEXT NOT NULL,
    video_url TEXT NOT NULL,
    verification_code VARCHAR(10) NOT NULL, -- Random code user speaks in video
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- 'pending', 'approved', 'rejected'
    rejection_reason TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    reviewed_at TIMESTAMP,
    reviewed_by UUID REFERENCES users(id)
);

-- Index for looking up by user
CREATE INDEX idx_verification_requests_user_id ON verification_requests(user_id);

-- Index for finding pending requests
CREATE INDEX idx_verification_requests_status ON verification_requests(status);

-- Only allow one pending request per user
CREATE UNIQUE INDEX idx_verification_requests_user_pending 
ON verification_requests(user_id) 
WHERE status = 'pending';



