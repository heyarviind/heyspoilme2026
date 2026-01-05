/**
 * Image compression and processing utilities
 */

export interface ImageProcessingOptions {
	maxWidth?: number;
	maxHeight?: number;
	quality?: number; // 0-1 for WebP quality
	maxFileSizeMB?: number;
}

const DEFAULT_OPTIONS: Required<ImageProcessingOptions> = {
	maxWidth: 1200,
	maxHeight: 1600,
	quality: 0.85,
	maxFileSizeMB: 10,
};

/**
 * Validates an image file before processing
 */
export function validateImage(file: File, maxSizeMB: number = 10): { valid: boolean; error?: string } {
	// Check if it's an image
	if (!file.type.startsWith('image/')) {
		return { valid: false, error: 'Please select an image file' };
	}

	// Check file size (before compression)
	const sizeMB = file.size / (1024 * 1024);
	if (sizeMB > maxSizeMB) {
		return { valid: false, error: `Image too large. Maximum size is ${maxSizeMB}MB` };
	}

	// Supported formats
	const supportedTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/webp', 'image/gif', 'image/heic', 'image/heif'];
	if (!supportedTypes.some(type => file.type.toLowerCase().includes(type.split('/')[1]))) {
		return { valid: false, error: 'Unsupported image format. Use JPEG, PNG, WebP, or GIF' };
	}

	return { valid: true };
}

/**
 * Loads an image file into an HTMLImageElement
 */
function loadImage(file: File): Promise<HTMLImageElement> {
	return new Promise((resolve, reject) => {
		const img = new Image();
		img.onload = () => resolve(img);
		img.onerror = () => reject(new Error('Failed to load image'));
		img.src = URL.createObjectURL(file);
	});
}

/**
 * Gets EXIF orientation from image file
 */
async function getExifOrientation(file: File): Promise<number> {
	try {
		const buffer = await file.slice(0, 65536).arrayBuffer();
		const view = new DataView(buffer);

		// Check for JPEG
		if (view.getUint16(0, false) !== 0xFFD8) return 1;

		let offset = 2;
		while (offset < view.byteLength) {
			const marker = view.getUint16(offset, false);
			offset += 2;

			if (marker === 0xFFE1) {
				// EXIF marker
				const length = view.getUint16(offset, false);
				offset += 2;

				// Check for "Exif"
				if (view.getUint32(offset, false) !== 0x45786966) return 1;
				offset += 6;

				// TIFF header
				const little = view.getUint16(offset, false) === 0x4949;
				offset += view.getUint32(offset + 4, little);

				const tags = view.getUint16(offset, little);
				offset += 2;

				for (let i = 0; i < tags; i++) {
					const tag = view.getUint16(offset + i * 12, little);
					if (tag === 0x0112) {
						return view.getUint16(offset + i * 12 + 8, little);
					}
				}
			} else if ((marker & 0xFF00) === 0xFF00) {
				offset += view.getUint16(offset, false);
			} else {
				break;
			}
		}
	} catch {
		// Ignore EXIF errors
	}
	return 1;
}

/**
 * Applies EXIF orientation to canvas context
 */
function applyOrientation(
	ctx: CanvasRenderingContext2D,
	width: number,
	height: number,
	orientation: number
): { width: number; height: number } {
	switch (orientation) {
		case 2:
			ctx.transform(-1, 0, 0, 1, width, 0);
			break;
		case 3:
			ctx.transform(-1, 0, 0, -1, width, height);
			break;
		case 4:
			ctx.transform(1, 0, 0, -1, 0, height);
			break;
		case 5:
			ctx.transform(0, 1, 1, 0, 0, 0);
			return { width: height, height: width };
		case 6:
			ctx.transform(0, 1, -1, 0, height, 0);
			return { width: height, height: width };
		case 7:
			ctx.transform(0, -1, -1, 0, height, width);
			return { width: height, height: width };
		case 8:
			ctx.transform(0, -1, 1, 0, 0, width);
			return { width: height, height: width };
	}
	return { width, height };
}

/**
 * Compresses and converts an image to WebP
 */
export async function compressImage(
	file: File,
	options: ImageProcessingOptions = {}
): Promise<{ blob: Blob; width: number; height: number }> {
	const opts = { ...DEFAULT_OPTIONS, ...options };

	// Load image
	const img = await loadImage(file);
	const originalWidth = img.naturalWidth;
	const originalHeight = img.naturalHeight;

	// Get EXIF orientation
	const orientation = await getExifOrientation(file);

	// Calculate new dimensions
	let { width, height } = calculateDimensions(
		originalWidth,
		originalHeight,
		opts.maxWidth,
		opts.maxHeight
	);

	// Swap dimensions for rotated images
	if (orientation >= 5 && orientation <= 8) {
		[width, height] = [height, width];
	}

	// Create canvas
	const canvas = document.createElement('canvas');
	canvas.width = width;
	canvas.height = height;

	const ctx = canvas.getContext('2d');
	if (!ctx) throw new Error('Could not get canvas context');

	// Apply orientation and get final dimensions
	const finalDims = applyOrientation(ctx, width, height, orientation);
	canvas.width = finalDims.width;
	canvas.height = finalDims.height;

	// Enable image smoothing for better quality
	ctx.imageSmoothingEnabled = true;
	ctx.imageSmoothingQuality = 'high';

	// Draw image
	ctx.drawImage(img, 0, 0, width, height);

	// Clean up object URL
	URL.revokeObjectURL(img.src);

	// Convert to WebP with quality adjustment
	let quality = opts.quality;
	let blob: Blob | null = null;

	// Try progressively lower quality if file is too large
	while (quality >= 0.3) {
		blob = await new Promise<Blob | null>((resolve) => {
			canvas.toBlob(resolve, 'image/webp', quality);
		});

		if (blob && blob.size <= 500 * 1024) {
			// Under 500KB is good
			break;
		}

		quality -= 0.1;
	}

	if (!blob) {
		throw new Error('Failed to compress image');
	}

	return {
		blob,
		width: finalDims.width,
		height: finalDims.height,
	};
}

/**
 * Calculates new dimensions while maintaining aspect ratio
 */
function calculateDimensions(
	originalWidth: number,
	originalHeight: number,
	maxWidth: number,
	maxHeight: number
): { width: number; height: number } {
	let width = originalWidth;
	let height = originalHeight;

	// Scale down if needed
	if (width > maxWidth) {
		height = Math.round((height * maxWidth) / width);
		width = maxWidth;
	}

	if (height > maxHeight) {
		width = Math.round((width * maxHeight) / height);
		height = maxHeight;
	}

	return { width, height };
}

/**
 * Generates a WebP filename from the original
 */
export function getWebPFilename(originalName: string): string {
	const baseName = originalName.replace(/\.[^/.]+$/, '');
	return `${baseName}.webp`;
}



