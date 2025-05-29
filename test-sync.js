// Simple test script for session synchronization
// Run with: node test-sync.js

const BACKEND_URL = "http://localhost:8080";

async function testSync() {
	console.log("🔄 Testing Backend Session Synchronization\n");

	// Test 1: Health check
	console.log("1. Testing health endpoint...");
	try {
		const response = await fetch(`${BACKEND_URL}/health`);
		const data = await response.json();
		console.log("   ✅ Health check:", data);
	} catch (error) {
		console.log("   ❌ Health check failed:", error.message);
	}

	// Test 2: Protected endpoint without auth
	console.log("\n2. Testing protected endpoint without auth...");
	try {
		const response = await fetch(`${BACKEND_URL}/api/acronyms`);
		if (response.status === 401) {
			console.log("   ✅ Correctly rejected unauthenticated request");
		} else {
			console.log("   ❌ Unexpected response:", response.status);
		}
	} catch (error) {
		console.log("   ❌ Error:", error.message);
	}

	// Test 3: Session invalidation endpoint
	console.log("\n3. Testing session invalidation endpoint...");
	try {
		const response = await fetch(`${BACKEND_URL}/api/session/invalidate`, {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify({ sessionToken: "test-token" }),
		});
		const data = await response.json();
		console.log("   ✅ Session invalidation:", data);
	} catch (error) {
		console.log("   ❌ Error:", error.message);
	}

	console.log("\n📋 Next steps:");
	console.log("   1. Start your frontend with better-auth");
	console.log("   2. Log in through the frontend");
	console.log(
		"   3. Use the frontend sync utility to test real authentication",
	);
	console.log("   4. Check browser cookies for better-auth.session_token");

	console.log("\n🔧 Frontend integration:");
	console.log('   import { BackendSync } from "@/lib/backend-sync";');
	console.log("   const status = await BackendSync.checkAuthStatus();");
}

testSync().catch(console.error);
