"use client";

import { useAuth } from "@/contexts/AuthContext";
import { redirect } from "next/navigation";

export default function AdminDashboard() {
  const { user, isAuthenticated } = useAuth();

  if (!isAuthenticated || !user) {
    redirect("/login");
  }

  // Check if user is admin
  if (user.roles && !user.roles.includes("admin")) {
    redirect("/dashboard");
  }

  return (
    <div>
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-white">Admin Dashboard</h1>
        <p className="text-gray-400 mt-2">Manage users, roles, and permissions</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        {/* Quick Stats */}
        <div className="bg-gray-800 rounded-lg p-6 border border-gray-700">
          <h3 className="text-gray-400 text-sm font-medium">Users</h3>
          <p className="text-2xl font-bold text-white mt-2">--</p>
          <p className="text-xs text-gray-500 mt-1">Loading...</p>
        </div>

        <div className="bg-gray-800 rounded-lg p-6 border border-gray-700">
          <h3 className="text-gray-400 text-sm font-medium">Roles</h3>
          <p className="text-2xl font-bold text-white mt-2">4</p>
          <p className="text-xs text-gray-500 mt-1">Default roles</p>
        </div>

        <div className="bg-gray-800 rounded-lg p-6 border border-gray-700">
          <h3 className="text-gray-400 text-sm font-medium">Permissions</h3>
          <p className="text-2xl font-bold text-white mt-2">16</p>
          <p className="text-xs text-gray-500 mt-1">Total permissions</p>
        </div>

        <div className="bg-gray-800 rounded-lg p-6 border border-gray-700">
          <h3 className="text-gray-400 text-sm font-medium">Admin Users</h3>
          <p className="text-2xl font-bold text-white mt-2">--</p>
          <p className="text-xs text-gray-500 mt-1">Loading...</p>
        </div>
      </div>

      <div className="mt-8 grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Quick Links */}
        <div className="bg-gray-800 rounded-lg p-6 border border-gray-700">
          <h2 className="text-lg font-semibold text-white mb-4">Quick Actions</h2>
          <ul className="space-y-2">
            <li>
              <a
                href="/admin/users"
                className="text-blue-400 hover:text-blue-300 flex items-center"
              >
                <span className="mr-2">→</span>Manage Users
              </a>
            </li>
            <li>
              <a
                href="/admin/roles"
                className="text-blue-400 hover:text-blue-300 flex items-center"
              >
                <span className="mr-2">→</span>Manage Roles
              </a>
            </li>
            <li>
              <a
                href="/admin/permissions"
                className="text-blue-400 hover:text-blue-300 flex items-center"
              >
                <span className="mr-2">→</span>View Permissions
              </a>
            </li>
            <li>
              <a
                href="/admin/assignments"
                className="text-blue-400 hover:text-blue-300 flex items-center"
              >
                <span className="mr-2">→</span>Role Assignments
              </a>
            </li>
          </ul>
        </div>

        {/* System Info */}
        <div className="bg-gray-800 rounded-lg p-6 border border-gray-700">
          <h2 className="text-lg font-semibold text-white mb-4">
            System Information
          </h2>
          <div className="space-y-2 text-sm">
            <div className="flex justify-between">
              <span className="text-gray-400">API Version:</span>
              <span className="text-white">v1</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-400">Database:</span>
              <span className="text-white">PostgreSQL 18.1</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-400">Cache:</span>
              <span className="text-white">Redis 8.4</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-400">Your Role:</span>
              <span className="text-green-400 font-medium">
                {user?.roles?.[0]?.charAt(0).toUpperCase() +
                  user?.roles?.[0]?.slice(1) || "Unknown"}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
