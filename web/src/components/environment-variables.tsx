"use client"

import * as React from "react"
import { IconPlus, IconTrash, IconEye, IconEyeOff } from "@tabler/icons-react"
import { Button } from "@/components/ui/button"
import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
} from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"

type EnvVariable = {
    id: string
    key: string
    value: string
    isVisible: boolean
}

export function EnvironmentVariables() {
    const [variables, setVariables] = React.useState<EnvVariable[]>([
        { id: "1", key: "DATABASE_URL", value: "postgresql://user:pass@localhost:5432/db", isVisible: false },
        { id: "2", key: "API_KEY", value: "sk_test_1234567890abcdef", isVisible: false },
    ])

    const [newKey, setNewKey] = React.useState("")
    const [newValue, setNewValue] = React.useState("")

    const addVariable = () => {
        if (!newKey.trim() || !newValue.trim()) return

        const newVar: EnvVariable = {
            id: Date.now().toString(),
            key: newKey.trim(),
            value: newValue.trim(),
            isVisible: false,
        }

        setVariables([...variables, newVar])
        setNewKey("")
        setNewValue("")
    }

    const removeVariable = (id: string) => {
        setVariables(variables.filter((v) => v.id !== id))
    }

    const toggleVisibility = (id: string) => {
        setVariables(
            variables.map((v) =>
                v.id === id ? { ...v, isVisible: !v.isVisible } : v
            )
        )
    }

    const updateVariable = (id: string, field: "key" | "value", value: string) => {
        setVariables(
            variables.map((v) =>
                v.id === id ? { ...v, [field]: value } : v
            )
        )
    }

    return (
        <Card>
            <CardHeader>
                <CardTitle>Environment Variables</CardTitle>
                <CardDescription>
                    Manage environment variables for your application. These will be available at runtime.
                </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
                {variables.map((variable) => (
                    <div key={variable.id} className="flex items-end gap-2">
                        <div className="flex-1 space-y-2">
                            <Label htmlFor={`key-${variable.id}`}>Key</Label>
                            <Input
                                id={`key-${variable.id}`}
                                value={variable.key}
                                onChange={(e) => updateVariable(variable.id, "key", e.target.value)}
                                placeholder="VARIABLE_NAME"
                                className="font-mono"
                            />
                        </div>
                        <div className="flex-1 space-y-2">
                            <Label htmlFor={`value-${variable.id}`}>Value</Label>
                            <div className="relative">
                                <Input
                                    id={`value-${variable.id}`}
                                    type={variable.isVisible ? "text" : "password"}
                                    value={variable.value}
                                    onChange={(e) => updateVariable(variable.id, "value", e.target.value)}
                                    placeholder="value"
                                    className="font-mono pr-10"
                                />
                                <Button
                                    type="button"
                                    variant="ghost"
                                    size="icon"
                                    className="absolute right-0 top-0 h-full"
                                    onClick={() => toggleVisibility(variable.id)}
                                >
                                    {variable.isVisible ? (
                                        <IconEyeOff className="h-4 w-4" />
                                    ) : (
                                        <IconEye className="h-4 w-4" />
                                    )}
                                </Button>
                            </div>
                        </div>
                        <Button
                            type="button"
                            variant="destructive"
                            size="icon"
                            onClick={() => removeVariable(variable.id)}
                        >
                            <IconTrash className="h-4 w-4" />
                        </Button>
                    </div>
                ))}

                <div className="border-t pt-4 space-y-4">
                    <h3 className="text-sm font-medium">Add New Variable</h3>
                    <div className="flex items-end gap-2">
                        <div className="flex-1 space-y-2">
                            <Label htmlFor="new-key">Key</Label>
                            <Input
                                id="new-key"
                                value={newKey}
                                onChange={(e) => setNewKey(e.target.value)}
                                placeholder="VARIABLE_NAME"
                                className="font-mono"
                            />
                        </div>
                        <div className="flex-1 space-y-2">
                            <Label htmlFor="new-value">Value</Label>
                            <Input
                                id="new-value"
                                value={newValue}
                                onChange={(e) => setNewValue(e.target.value)}
                                placeholder="value"
                                className="font-mono"
                            />
                        </div>
                        <Button
                            type="button"
                            onClick={addVariable}
                            disabled={!newKey.trim() || !newValue.trim()}
                        >
                            <IconPlus className="h-4 w-4 mr-2" />
                            Add
                        </Button>
                    </div>
                </div>

                <div className="flex justify-end pt-4">
                    <Button>Save Changes</Button>
                </div>
            </CardContent>
        </Card>
    )
}
