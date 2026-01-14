import Link from "next/link";

export default function Home() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-gray-900 to-gray-800">
      <div className="text-center space-y-8 p-8">
        <div>
          <h1 className="text-6xl font-extrabold text-white mb-4">
            Venio
          </h1>
          <p className="text-xl text-gray-400 mb-8">
            User Management System
          </p>
        </div>

        <div className="flex gap-4 justify-center">
          <Link
            href="/login"
            className="px-8 py-3 text-lg font-medium text-white bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors"
          >
            Sign In
          </Link>
          <Link
            href="/register"
            className="px-8 py-3 text-lg font-medium text-white bg-gray-700 hover:bg-gray-600 rounded-lg transition-colors"
          >
            Register
          </Link>
        </div>

        <div className="mt-12 text-sm text-gray-500">
          <p>A modern authentication system built with Go and Next.js</p>
        </div>
      </div>
    </div>
  );
}
