# AI Integration Documentation

This document explains how to use the AI integration features in the blog API.

## Overview

The AI integration provides two main features:
1. **Content Suggestions**: Get AI-generated suggestions for improving existing blog content
2. **Content Ideas**: Generate new blog post ideas based on keywords

## Prerequisites

1. **AI Service Running**: Ensure the AI service is running on `http://localhost:8001`
2. **Environment Configuration**: Set `AI_SERVICE_URL=http://localhost:8001` in your `.env` file
3. **Authentication**: All AI endpoints require user authentication

## API Endpoints

### 1. Generate Content Suggestions

**Endpoint**: `POST /api/ai/suggestions`

**Authentication**: Required (Bearer Token)

**Request Body**:
```json
{
  "keywords": ["technology", "programming", "golang"],
  "tone": "professional",
  "type": "improvement",
  "blog_content": "Your existing blog content here..."
}
```

### 2. Save AI Suggestion

**Endpoint**: `POST /api/ai/save`

**Authentication**: Required (Bearer Token)

**Request Body**:
```json
{
  "input_topic": "Advanced Go Programming",
  "keywords": ["golang", "programming", "concurrency"],
  "tone": "professional",
  "suggestions": [
    "Title: Advanced Go Programming Techniques",
    "Target Audience: Software developers and engineers",
    "• Understanding Go concurrency patterns",
    "• Best practices for Go performance optimization"
  ]
}
```

### 3. Get User's AI Suggestions

**Endpoint**: `GET /api/ai/suggestions?page=1&limit=10`

**Authentication**: Required (Bearer Token)

**Query Parameters**:
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 10, max: 100)

### 4. Get AI Suggestions by Status

**Endpoint**: `GET /api/ai/suggestions/status/{status}?page=1&limit=10`

**Authentication**: Required (Bearer Token)

**Path Parameters**:
- `status`: One of "saved", "converted-to-draft", "discarded"

**Query Parameters**:
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 10, max: 100)

### 5. Convert Suggestion to Draft

**Endpoint**: `POST /api/ai/suggestions/{id}/convert-to-draft`

**Authentication**: Required (Bearer Token)

**Path Parameters**:
- `id`: AI suggestion ID

### 6. Delete AI Suggestion

**Endpoint**: `DELETE /api/ai/suggestions/{id}`

**Authentication**: Required (Bearer Token)

**Path Parameters**:
- `id`: AI suggestion ID

**Response**:
```json
{
  "success": true,
  "message": "AI suggestions generated successfully",
  "data": {
    "suggestions": [
      "Title: Advanced Go Programming Techniques",
      "Target Audience: Software developers and engineers",
      "• Understanding Go concurrency patterns",
      "• Best practices for Go performance optimization",
      "• Modern Go development workflows"
    ],
    "type": "improvement",
    "message": "AI suggestions generated successfully"
  }
}
```

### 7. Generate Content Ideas

**Endpoint**: `POST /api/ai/ideas`

**Authentication**: Required (Bearer Token)

**Request Body**:
```json
{
  "keywords": ["artificial intelligence", "machine learning"],
  "tone": "educational"
}
```

**Response**:
```json
{
  "success": true,
  "message": "Content ideas generated successfully",
  "data": {
    "suggestions": [
      "Title: Introduction to Machine Learning for Beginners",
      "Target Audience: Students and professionals new to AI",
      "• Understanding the basics of machine learning",
      "• Common algorithms and their applications",
      "• Getting started with AI development"
    ],
    "type": "ideas",
    "message": "Content ideas generated successfully"
  }
}
```

## Usage Examples

### Using cURL

1. **Get Content Suggestions**:
```bash
curl -X POST http://localhost:8080/api/ai/suggestions \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "keywords": ["web development", "react"],
    "tone": "casual",
    "blog_content": "Your blog content here..."
  }'
```

2. **Save AI Suggestion**:
```bash
curl -X POST http://localhost:8080/api/ai/save \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "input_topic": "React Development",
    "keywords": ["react", "javascript", "frontend"],
    "tone": "professional",
    "suggestions": [
      "Title: Modern React Development Patterns",
      "Target Audience: Frontend developers",
      "• Understanding React hooks and context",
      "• Performance optimization techniques"
    ]
  }'
```

3. **Get User's Suggestions**:
```bash
curl -X GET "http://localhost:8080/api/ai/suggestions?page=1&limit=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

4. **Convert Suggestion to Draft**:
```bash
curl -X POST http://localhost:8080/api/ai/suggestions/SUGGESTION_ID/convert-to-draft \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

5. **Generate Content Ideas**:
```bash
curl -X POST http://localhost:8080/api/ai/ideas \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "keywords": ["blockchain", "cryptocurrency"],
    "tone": "professional"
  }'
```

### Using JavaScript/Fetch

```javascript
// Get content suggestions
const response = await fetch('http://localhost:8080/api/ai/suggestions', {
  method: 'POST',
  headers: {
    'Authorization': 'Bearer YOUR_JWT_TOKEN',
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    keywords: ['technology', 'programming'],
    tone: 'professional',
    blog_content: 'Your existing blog content...'
  })
});

const data = await response.json();
console.log(data.data.suggestions);
```

## Error Handling

The API returns standardized error responses:

```json
{
  "success": false,
  "message": "Error description"
}
```

Common error scenarios:
- **400 Bad Request**: Invalid input format or missing required fields
- **401 Unauthorized**: Missing or invalid authentication token
- **500 Internal Server Error**: AI service unavailable or configuration issues

## Configuration

### Environment Variables

Add the following to your `.env` file:
```
AI_SERVICE_URL=http://localhost:8001
```

### AI Service Requirements

The AI service should:
- Run on the configured port (default: 8001)
- Accept POST requests to `/generate`
- Expect request body with `keywords` and `tone` fields
- Return JSON response with `title`, `audience`, and `headlines` fields

## User Experience Flow

### Complete AI Suggestion Workflow

1. **Generate Suggestion**:
   - User enters topic/keywords
   - System calls AI service to generate suggestions
   - User receives AI-generated content ideas

2. **Save or Discard**:
   - User can save the suggestion to their personal collection
   - User can discard the suggestion (not stored)
   - Saved suggestions are stored with status "saved"

3. **Manage Suggestions**:
   - View all saved suggestions
   - Filter by status (saved, converted-to-draft, discarded)
   - Delete unwanted suggestions

4. **Convert to Draft**:
   - User can convert a saved suggestion to a blog draft
   - Creates a new blog post with draft status
   - Suggestion status changes to "converted-to-draft"

5. **Edit and Publish**:
   - User can edit the draft blog post
   - Add content, modify title, add tags
   - Publish when ready

### Example User Journey

```
User enters topic "Go Programming" 
    ↓
AI generates suggestions
    ↓
User sees suggestions with buttons:
    ✅ "Save as draft" → stored in DB
    🗑 "Discard" → nothing saved
    ↓
If saved, user can:
    - View in their suggestions list
    - Convert to blog draft
    - Edit the draft
    - Publish the blog post
```

## Security Considerations

1. **Authentication**: All AI endpoints require valid JWT tokens
2. **Input Validation**: All inputs are validated before processing
3. **Error Handling**: Sensitive error details are not exposed to clients
4. **Rate Limiting**: Consider implementing rate limiting for AI endpoints
5. **Data Privacy**: Users can only access their own suggestions

## Troubleshooting

1. **AI Service Connection Issues**:
   - Verify AI service is running on the correct port
   - Check `AI_SERVICE_URL` environment variable
   - Ensure network connectivity between services

2. **Authentication Issues**:
   - Verify JWT token is valid and not expired
   - Check token format (Bearer token)

3. **Response Parsing Issues**:
   - Ensure AI service returns expected JSON format
   - Check response structure matches expected fields
