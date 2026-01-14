"use client";

import { useAuth } from "@/contexts/AuthContext";
import { redirect } from "next/navigation";
import { useState, useEffect } from "react";

interface Permission {
  id: string;
  name: string;
  description: string;
  resource: string;
  action: string;
  created_at: string;
}

export default function PermissionsPage() {
  const { user, isAuthenticated, token } = useAuth();
  const [permissions, setPermissions] = useState<Permission[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [filter, setFilter] = useState("all");

  if (!isAuthenticated || !user) {
    redirect("/login");
  }

  if (user.roles && !user.roles.includes("admin")) {
    redirect("/dashboard");
  }

  useEffect(() => {
    fetchPermissions();
  }, [token]);

  const fetchPermissions = async () => {
    try {
      setLoading(true);
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/api/v1/admin/permissions`,
        {
          headers: {
            "Authorization": `Bearer ${token}`,
            "Content-Type": "application/json",
          },
        }
      );

      if (!response.ok) {
        throw new Error(`Failed to fetch permissions: ${response.statusText}`);
      }

      const data = await response.json();
      setPermissions(data.permissions || []);
      setError(null);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Failed to fetch permissions"
      );
    } finally {
      setLoading(false);
    }
  };

  const resources = Array.from(
    new Set(permissions.map((p) => p.resource))
  ).sort();

  const filteredPermissions =
    filter === "all"
      ? permissions
      : permissions.filter((p) => p.resource === filter);

  return (
    <div>
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-white">Permissions</h1>
        <p className="text-gray-400 mt-2">
          View all available permissions in the system
        </p>
      </div>

      {error && (
        <div className="bg-red-900 border border-red-700 text-red-100 px-4 py-3 rounded-lg mb-4">
          {error}
        </div>
      )}

      {!loading && (
        <div className="mb-4">
          <label className="block text-sm font-medium text-gray-300 mb-2">
            Filter by Resource
          </label>
          <select
            value={filter}
            onChange={(e) => setFilter(e.target.value)}
            className="bg-gray-700 text-white px-4 py-2 rounded-lg border border-gray-600 focus:border-blue-500 focus:outline-none"
          >
            <option value="all">All Resources</option>
            {resources.map((resource) => (
              <option key={resource} value={resource}>
                {resource.charAt(0).toUpperCase() + resource.slice(1)}
              </option>
            ))}
          </select>
        </div>
      )}

      {loading ? (
        <div className="bg-gray-800 rounded-lg p-8 text-center border border-gray-700">
          <p className="text-gray-400">Loading permissions...</p>
        </div>
      ) : filteredPermissions.length === 0 ? (
        <div className="bg-gray-800 rounded-lg p-8 text-center border border-gray-700">
          <p className="text-gray-400">No permissions found</p>
        </div>
      ) : (
        <div className="bg-gray-800 border border-gray-700 rounded-lg overflow-hidden">
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4 p-4 bg-gray-700 font-semibold text-white text-sm">
            <div>Resource</div>
            <div>Action</div>
            <div>Name</div>
            <div>Description</div>
          </div>
          <div className="divide-y divide-gray-700">
            {filteredPermissions.map((permission) => (
              <div
                key={permission.id}
                className="grid grid-cols-1 md:grid-cols-4 gap-4 p-4 hover:bg-gray-700 transition"
              >
                <div className="font-medium text-blue-400">
                  {permission.resource}
                </div>
                <div className="text-gray-300">{permission.action}</div>
                <div className="text-white">{permission.name}</div>
                <div className="text-gray-400 text-sm">
                  {permission.description}
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
